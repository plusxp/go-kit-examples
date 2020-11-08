package client

import (
	stdendpoint "github.com/go-kit/kit/endpoint"
	stdopentracing "github.com/opentracing/opentracing-go"
	"io"
	config2 "new_addsvc/config"
	endpoint2 "new_addsvc/pkg/endpoint"
	service2 "new_addsvc/pkg/service"
	transport2 "new_addsvc/pkg/transport"
	"time"

	consulapi "github.com/hashicorp/consul/api"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
)

// New returns a service that's load-balanced over instances of new_addsvc found
// in the provided Consul server. The mechanism of looking up new_addsvc
// instances in Consul is hard-coded into thient.
// client从consul获取实例地址
func New(consulAddr string, logger log.Logger) (service2.Service, error) {
	apiClient, err := consulapi.NewClient(&consulapi.Config{
		Address: consulAddr,
	})
	if err != nil {
		return nil, err
	}

	// As the implementer of new_addsvc, we declare and enforce these
	// parameters for all of the new_addsvc consumers.
	var (
		consulService = config2.SvcName
		consulTags    = []string{"gokit_svc"}
		passingOnly   = true // 只获取健康的实例地址
		retryMax      = 3
		retryTimeout  = 500 * time.Millisecond
	)

	var (
		sdclient  = consul.NewClient(apiClient)
		instancer = consul.NewInstancer(sdclient, logger, consulService, consulTags, passingOnly)
		/*
			client得到的对象还是endpoint
		*/
		endpoints endpoint2.AddSvcEndpoints
	)

	var tracer stdopentracing.Tracer
	tracer = stdopentracing.GlobalTracer()

	{
		// 在client，每个endpoint又依次封装了服务发现、负载均衡、重试，还可以加断路器，限速等
		// 每个endpoint单独封装，可以非常细粒度的为接口安装基础设施（比如某些接口的限速配置与其他接口并不相同）
		factory := factoryFor(tracer, log.NewNopLogger(), endpoint2.MakeSumEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.SumEndpoint = retry
	}
	{
		factory := factoryFor(tracer, log.NewNopLogger(), endpoint2.MakeConcatEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.ConcatEndpoint = retry
	}

	return endpoints, nil
}

type MakeEndpoint func(service2.Service) stdendpoint.Endpoint

func factoryFor(otTracer stdopentracing.Tracer, logger log.Logger, makeEndpoint MakeEndpoint) sd.Factory {
	return func(instance string) (stdendpoint.Endpoint, io.Closer, error) {
		svc, err := transport2.MakeClientEndpoints(instance, otTracer, logger)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(svc), nil, nil
	}
}
