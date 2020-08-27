// Package client provides a profilesvc client based on a predefined Consul
// service name and relevant tags. Users must only provide the address of a
// Consul server.
package client

import (
	stdendpoint "github.com/go-kit/kit/endpoint"
	stdopentracing "github.com/opentracing/opentracing-go"
	"go-kit-examples/new_addsvc/pkg/service"
	"go-kit-examples/new_addsvc/pkg/transport"
	"io"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	"go-kit-examples/new_addsvc/pkg/endpoint"
)

// New returns a service that's load-balanced over instances of profilesvc found
// in the provided Consul server. The mechanism of looking up profilesvc
// instances in Consul is hard-coded into thient.
// client从consul获取实例地址
func New(consulAddr string, logger log.Logger) (service.Service, error) {
	apiclient, err := consulapi.NewClient(&consulapi.Config{
		Address: consulAddr,
	})
	if err != nil {
		return nil, err
	}

	// As the implementer of profilesvc, we declare and enforce these
	// parameters for all of the profilesvc consumers.
	var (
		consulService = "addsvc"
		consulTags    = []string{"prod"}
		passingOnly   = true
		retryMax      = 3
		retryTimeout  = 500 * time.Millisecond
	)

	var (
		sdclient  = consul.NewClient(apiclient)
		instancer = consul.NewInstancer(sdclient, logger, consulService, consulTags, passingOnly)
		/*
			client得到的对象还是endpoint，
		*/
		endpoints endpoint.AddSvcEndpoints
	)

	var tracer stdopentracing.Tracer
	tracer = stdopentracing.GlobalTracer()

	{
		// 在client，每个endpoint又依次封装了服务发现、负载均衡、重试，还可以加断路器，限速等
		// 每个endpoint单独封装，可以非常细粒度的为接口安装基础设施（比如某些接口的限速配置与其他接口并不相同）
		factory := factoryFor(tracer, log.NewNopLogger(), endpoint.MakeSumEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.SumEndpoint = retry
	}
	{
		factory := factoryFor(tracer, log.NewNopLogger(), endpoint.MakeConcatEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.ConcatEndpoint = retry
	}

	return endpoints, nil
}

type MakeEndpoint func(service.Service) stdendpoint.Endpoint

func factoryFor(otTracer stdopentracing.Tracer, logger log.Logger, makeEndpoint MakeEndpoint) sd.Factory {
	return func(instance string) (stdendpoint.Endpoint, io.Closer, error) {
		svc, err := transport.MakeClientEndpoints(instance, otTracer, logger)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(svc), nil, nil
	}
}
