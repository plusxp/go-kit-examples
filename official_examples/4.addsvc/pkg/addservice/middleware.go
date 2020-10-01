package addservice

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// LoggingMiddleware takes a loggermw as a dependency
// and returns a service Middleware.
func LoggingMiddleware(loggermw log.Logger, ints, chars metrics.Counter) Middleware {
	return func(next Service) Service {
		instrumw := instrumentingMiddleware{
			ints:  ints,
			chars: chars,
			next:  next,
		}
		return unifyMiddleware{loggermw, instrumw, next}
	}
}

type instrumentingMiddleware struct {
	ints  metrics.Counter
	chars metrics.Counter
	next  Service
}

type unifyMiddleware struct {
	// 通过命名规范代码， 不同功能的对象通过嵌套struct添加
	loggermw log.Logger
	instrumw instrumentingMiddleware
	next     Service
}

func (mw unifyMiddleware) Sum(ctx context.Context, a, b int) (v int, err error) {
	defer func() {
		mw.loggermw.Log("method", "Sum", "a", a, "b", b, "v", v, "err", err)
	}()
	v, err = mw.next.Sum(ctx, a, b)
	mw.instrumw.ints.Add(float64(v))
	return v, err
}

func (mw unifyMiddleware) Concat(ctx context.Context, a, b string) (v string, err error) {
	defer func() {
		mw.loggermw.Log("method", "Concat", "a", a, "b", b, "v", v, "err", err)
	}()
	return mw.next.Concat(ctx, a, b)
}

// InstrumentingMiddleware returns a service middleware that instruments
// the number of integers summed and characters concatenated over the lifetime of
// the service.
func InstrumentingMiddleware(ints, chars metrics.Counter) Middleware {
	return func(next Service) Service {
		return instrumentingMiddleware{
			ints:  ints,
			chars: chars,
			next:  next,
		}
	}
}

func (mw instrumentingMiddleware) Sum(ctx context.Context, a, b int) (int, error) {
	v, err := mw.next.Sum(ctx, a, b)
	mw.ints.Add(float64(v))
	return v, err
}

func (mw instrumentingMiddleware) Concat(ctx context.Context, a, b string) (string, error) {
	v, err := mw.next.Concat(ctx, a, b)
	mw.chars.Add(float64(len(v)))
	return v, err
}
