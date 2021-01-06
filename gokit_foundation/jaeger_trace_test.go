package gokit_foundation

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"go-util/_util"
	"testing"
	"time"
)

func TestInitTracer(t *testing.T) {
	Svc := "TestInitTracer"
	closer, err := InitTracer(Svc)
	_util.PanicIfErr(err, nil)

	defer closer.Close()
	tagOpt := opentracing.Tag{
		Key:   "type",
		Value: "GET",
	}

	// 假设正在处理一个请求
	// 首先对这个请求进行身份验证
	rootSpan := opentracing.StartSpan("/Op-auth", tagOpt)

	defer rootSpan.Finish()
	// 无论成功失败，记录一下结果
	rootSpan.LogFields(log.String("result", "pass"), log.Int("uid", 1000123))

	// 可选，添加一个rpc client tag，表示这个请求即将发起一个RPC调用
	//（添加tag的方式有多种，这只是其中一种，但这里添加的tag的key是OpenTracing的规范）
	ext.SpanKindRPCClient.Set(rootSpan)

	// 然后调用了一个RPC接口getUserOrderList，获取它的订单列表, 下面是RPC接口的内部逻辑
	// 这里就开启一个child span
	childSpan := opentracing.StartSpan("/Op-getUserOrderList", opentracing.ChildOf(rootSpan.Context()))
	ext.SpanKindRPCServer.Set(childSpan)

	defer childSpan.Finish()
	childSpan.LogFields(log.Int("orderCount", 22))

	time.Sleep(time.Second) // 模拟RPC耗时
}
