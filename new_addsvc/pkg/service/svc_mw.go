package service

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

type Middleware func(Service) Service

// mw
func UnifyMiddleware(loggermw log.Logger, ints, chars metrics.Counter) Middleware {
	return func(next Service) Service {
		instrumw := instrumentingMiddleware{
			ints:  ints,
			chars: chars,
			next:  next,
		}
		return unifyMiddleware{loggermw, instrumw, next}
	}
}

// 监控mw
type instrumentingMiddleware struct {
	ints  metrics.Counter
	chars metrics.Counter
	next  Service
}

// mw实体
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
