package grpc

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gokit_foundation"
	"hello/pb/gen-go/pbcommon"
	"testing"
)

func TestNew(t *testing.T) {
	lgr := gokit_foundation.NewLogger(nil)

	svc := "UpstreamSvcOfHelloSvc"
	c := MustNew(lgr, svc)
	defer c.Close()

	// 出于演示效果，这里创建一个root span，
	rootSpan := opentracing.GlobalTracer().StartSpan("TestNew")
	ctx := opentracing.ContextWithSpan(context.TODO(), rootSpan)
	defer rootSpan.Finish()
	rootSpan.LogFields(log.Bool("test", true))

	// 在日常开发中如果是频繁使用同一个RPC client，则不需要使用完立即close，而是完全不需要时再close，频繁创建连接也是开销
	// 也可以使用sync.Pool封装client，调用者不再关心close问题
	reply, err := c.SayHi(ctx, "Jack Ma")
	if err != pbcommon.R_OK {
		t.Error(err)
	}
	lgr.Log("rsp", reply)
}
