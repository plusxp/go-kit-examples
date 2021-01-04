package gokit_foundation

import (
	opentracinggo "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"go-util/_util"
	"testing"
)

func TestInitTracer(t *testing.T) {
	Svc := "TestInitTracer"
	tracer, closer, err := InitTracer(Svc)
	_util.PanicIfErr(err, nil)

	defer closer.Close()

	tagOpt := opentracinggo.Tag{
		Key:   "type",
		Value: "GET",
	}
	// 假设正在处理一个请求
	// 首先对这个请求进行身份验证
	rootSpan := tracer.StartSpan("Op-auth", tagOpt)
	defer rootSpan.Finish()
	// 无论成功失败，记录一下结果
	rootSpan.LogFields(log.String("result", "pass"), log.Int("uid", 1000123))

	// 然后可以获取它的订单列表
	childSpan := tracer.StartSpan("Op-getUserOrderList", opentracinggo.ChildOf(rootSpan.Context()))
	defer childSpan.Finish()
	childSpan.LogFields(log.Int("orderCount", 22))
}
