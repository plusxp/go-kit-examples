package addendpoint

import (
	"context"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"go-kit-examples/new-to-gokit/5.addsvc_new/pkg/addservice"
	"golang.org/x/time/rate"
	"time"
)

// 这一层封装service的每个方法，添加的限速、断路器、跟踪器
func NewEP(sw *addservice.SvcWrapper, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) EpMap {
	var epm = make(EpMap)
	sw.Range(func(method string, h addservice.SvcHandler) {
		ep := makeEP(h)
		ep = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(ep)
		ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(ep)
		ep = opentracing.TraceServer(otTracer, method)(ep)
		if zipkinTracer != nil {
			ep = zipkin.TraceEndpoint(zipkinTracer, method)(ep)
		}
		ep = LoggingMiddleware(log.With(logger, "method", method))(ep)
		ep = InstrumentingMiddleware(duration.With("method", method))(ep)
		epm[method] = ep
	})
	return epm
}

func makeEP(h addservice.SvcHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return h(ctx, request)
	}
}
