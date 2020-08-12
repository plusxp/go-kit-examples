package main

import (
	"fmt"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"net/url"
	"strings"
	"time"
)

// Set some parameters for our client.
var (
	qps         = 100                    // beyond which we will return an error
	maxAttempts = 3                      // per request, before giving up
	maxTime     = 250 * time.Millisecond // wallclock time, before giving up
)

// 从服务发现获取多个 可连接地址
// 改造存在的问题：每次接口执行都获取一个相同列表，会导致永远调用第一个地址
func getInstancesFromServerDiscovery() string {
	return "localhost:8081,localhost:8080"
}

func (svc *stringService) getRPCEndpoint(route rpcRoute) endpoint.Endpoint {
	instances := getInstancesFromServerDiscovery()
	var (
		instanceList = split(instances)
		endpointer   sd.FixedEndpointer
	)
	// 构建多个endpoint（嵌入断路器，限速器）
	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = makeRPCEndpoint(route, instance)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		endpointer = append(endpointer, e)
	}

	// Now, build a single, retrying, load-balancing endpoint out of all of
	// those individual endpoints.
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)
	return retry
}

func makeRPCEndpoint(route rpcRoute, instance string) endpoint.Endpoint {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}
	if _, ok := rpcRoutes[route]; !ok {
		panic(fmt.Sprintf("route:%s not found", route))
	}
	u.Path = string(route)
	return httptransport.NewClient(
		"GET",
		u,
		encodeRequest,
		decodeUppercaseResponse,
	).Endpoint()
}

type rpcRoute string

const (
	route_uppercase rpcRoute = "/uppercase"
)

var rpcRoutes = map[rpcRoute]struct{}{
	route_uppercase: {},
}
