package service

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"go-util/_util"
	"gokit_foundation"
	"hello/config"
	"hello/db"
	"hello/pkg/crontask"
)

func initFirstly(logger *gokit_foundation.Logger) {
	db.Init(logger) // 先启动db
	var err error
	tracerCloser, err = gokit_foundation.InitTracer(config.SvcName)
	tracer = opentracing.GlobalTracer()
	_util.PanicIfErr(err, nil, fmt.Sprintf("InitTracer err %v", err))
	crontask.Init()
}

func onClose() {
	crontask.Stop()
	tracerCloser.Close()
	db.Close() // 最后停止db
}
