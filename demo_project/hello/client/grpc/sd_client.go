package grpc

import (
	"fmt"
	stdendpoint "github.com/go-kit/kit/endpoint"
	consulapi "github.com/hashicorp/consul/api"
	stdopentracing "github.com/opentracing/opentracing-go"
	"go-util/_str"
	"gokit_foundation"
	"hello/config"
	"hello/pkg/endpoint"
	grpctransport "hello/pkg/grpc"
	"hello/pkg/service"
	"io"
	"reflect"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
)

/*
	带服务发现功能的client
	SD: service discovery
*/

// 带SD功能的client不需要开发者管理底层conn，交由go-kit管理
var svcSdClient service.HelloService

func MustNewClientWithSd(logger *gokit_foundation.Logger) service.HelloService {
	if svcSdClient == nil {
		svcSdClient = newHelloClientWithSd("", logger)
	}
	return svcSdClient
}

// client从consul获取实例地址
func newHelloClientWithSd(consulAddr string, logger *gokit_foundation.Logger) service.HelloService {
	_str.SetDefault(&consulAddr, consulAddr, ":8500")
	// 这一步并不会尝试连接consul，仅做连接配置检查
	apiClient, err := consulapi.NewClient(&consulapi.Config{
		Address: consulAddr,
	})
	// panics if err is not nil
	logger.PanicIfErr(err, nil, fmt.Sprintf("consulapi.NewClient err:%v", err))

	/* 统一的mw的参数 */
	var (
		consulService = config.SvcName
		consulTags    = []string{"gokit_svc"}
		passingOnly   = true // 只获取健康的实例地址
		retryMax      = 3
		retryTimeout  = time.Second
	)

	var (
		sdclient  = consul.NewClient(apiClient)
		instancer = consul.NewInstancer(sdclient, logger, consulService, consulTags, passingOnly)
		endpoints endpoint.Endpoints
	)

	var tracer stdopentracing.Tracer
	tracer = stdopentracing.GlobalTracer()

	setEndpointWithSD := func(makeEP MakeEndpoint, method string) {
		factory := factoryFor(tracer, log.NewNopLogger(), makeEP)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		ep := reflect.ValueOf(&endpoints).Elem().FieldByName(method + "Endpoint")
		ep.Set(reflect.ValueOf(retry))
	}

	setEndpointWithSD(endpoint.MakeSayHiEndpoint, "SayHi")
	setEndpointWithSD(endpoint.MakeMakeADateEndpoint, "MakeADate")
	setEndpointWithSD(endpoint.MakeUpdateUserInfoEndpoint, "UpdateUserInfo")

	return endpoints
}

type MakeEndpoint func(service.HelloService) stdendpoint.Endpoint

func factoryFor(otTracer stdopentracing.Tracer, logger log.Logger, makeEndpoint MakeEndpoint) sd.Factory {
	return func(instance string) (stdendpoint.Endpoint, io.Closer, error) {
		svc, err := grpctransport.MakeClientEndpoints(instance, otTracer, logger)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(svc), nil, nil
	}
}
