package endpoint

import (
	"context"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	"go-kit-examples/new_addsvc/pb/gen-go/resultcode"
	"go-kit-examples/new_addsvc/pkg/service"
	"golang.org/x/time/rate"
	"time"
)

// eps 内包含的ep对应service每个接口，而且必须一致
// ep非常重要，在RPC调用时，client调用的也是ep，对调用者隐藏了transport层
type AddSvcEndpoints struct {
	SumEndpoint    endpoint.Endpoint
	ConcatEndpoint endpoint.Endpoint
}

// 将一个Service对象转为Endpoints对象
func New(svc service.Service, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer) AddSvcEndpoints {
	var sumEndpoint endpoint.Endpoint
	// 使用洋葱模式封装endpoint
	{
		sumEndpoint = MakeSumEndpoint(svc)

		sumEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(sumEndpoint)
		sumEndpoint = opentracing.TraceServer(otTracer, "Sum")(sumEndpoint)
		sumEndpoint = InstrumentingMiddleware(duration.With("method", "Sum"))(sumEndpoint)
	}

	var concatEndpoint endpoint.Endpoint
	{
		concatEndpoint = MakeConcatEndpoint(svc)

		concatEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Limit(1), 100))(concatEndpoint)
		concatEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(concatEndpoint)
		concatEndpoint = opentracing.TraceServer(otTracer, "Concat")(concatEndpoint)
		concatEndpoint = InstrumentingMiddleware(duration.With("method", "Concat"))(concatEndpoint)
	}
	return AddSvcEndpoints{
		SumEndpoint:    sumEndpoint,
		ConcatEndpoint: concatEndpoint,
	}
}

// 针对接口：Sum 的转换方法
func MakeSumEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*SumRequest)
		v, err := s.Sum(ctx, req.A, req.B)
		// err 在logger mw打印出来 以debug
		// 谨慎处理endpoint层返回的err，这个err会被各种中间件捕获，可能会产生一些影响
		// 比如断路器就是安装在endpoint层，它会根据返回err次数来判断什么时候打开开关
		// 所以，业务err应该封装在SumResponse内，如果是你在service层能够读取到接口调用的系统压力过大
		// 这个时候可以通过endpoint层返出去
		// 一般情况，这里都返回err
		return &SumResponse{RetCode: errToRetCode(err), V: v}, nil
	}
}

// 针对接口：Concat 的转换方法
func MakeConcatEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*ConcatRequest)
		p, err := s.Concat(ctx, req.A, req.B)
		return &ConcatResponse{RetCode: errToRetCode(err), V: p}, nil
	}
}

// 统一处理err
func errToRetCode(err error) resultcode.RESULT_CODE {
	switch err {
	case service.ErrIntOverflow, service.ErrMaxSizeExceeded, service.ErrTwoZeroes:
		return resultcode.RESULT_CODE_RET_INVALID_ARGS
	default:
		return resultcode.RESULT_CODE_RET_UNKNOWN_ERR
	}
}
