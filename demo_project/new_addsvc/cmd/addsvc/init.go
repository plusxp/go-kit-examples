package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/leigg-go/go-util/_redis"
	stdopentracing "github.com/opentracing/opentracing-go"
	"go-util/_util"
	"gokit_foundation"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"new_addsvc/config"
	"new_addsvc/internal"
	"new_addsvc/pb/gen-go/addsvcpb"
	"new_addsvc/pkg/crontask"
	"new_addsvc/pkg/endpoint"
	"new_addsvc/pkg/service"
	"new_addsvc/pkg/transport"
)

func beforeStart() {
	_redis.MustInitDef(config.GetRedisConf())
}

func afterStart(grpcHost string, grpcPort int) {
	crontask.Init()
	gokit_foundation.MustRegisterSvc(config.SvcName, grpcHost, grpcPort, []string{"test"}) // 后上线
}

func beforeStop() {
	gokit_foundation.ConsulDeregister() // 先下线
	crontask.Stop()
}

// 最后停止一些基础设施，一般是db相关、mq
func afterStop() {
	_redis.Close()
}

func startHttpSrv(mainCtx context.Context, httpSrvAddr string) <-chan struct{} {
	closed := make(chan struct{})
	httpSrv := &http.Server{}
	stop := func() {
		defer func() { closed <- struct{}{} }()
		<-mainCtx.Done()
		err := httpSrv.Shutdown(context.TODO())
		//if err != nil {}
		logger.Log(fmt.Sprintf("httpSrv(%s)", httpSrvAddr), "exited", "err", err)
	}
	go func() {
		logger.Log("startHttpSrv", httpSrvAddr)

		httpLis, err := net.Listen("tcp", httpSrvAddr)
		_util.PanicIfErr(err, nil)

		go stop()
		// default use http.DefaultServeMux as handler
		err = httpSrv.Serve(httpLis)
		_util.PanicIfErr(err, []error{http.ErrServerClosed})
	}()

	return closed
}

func startgRPCSrv(mainCtx context.Context, grpcSrvAddr string) <-chan struct{} {
	closed := make(chan struct{})
	gRPCSrv := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	stop := func() {
		defer func() { closed <- struct{}{} }()
		<-mainCtx.Done()
		gRPCSrv.GracefulStop()
		logger.Log(fmt.Sprintf("gRPCSrv(%s)", grpcSrvAddr), "exited")
	}
	go func() {
		logger.Log("startgRPCSrv", grpcSrvAddr)
		gRPCLis, err := net.Listen("tcp", grpcSrvAddr)
		_util.PanicIfErr(err, nil)

		addSrv := NewGokitSrv(logger)
		addsvcpb.RegisterAddServer(gRPCSrv, addSrv)

		gokit_foundation.RegistergRPCHealthSrv(gRPCSrv) // 这里注册healthSrv
		go stop()
		err = gRPCSrv.Serve(gRPCLis)
		_util.PanicIfErr(err, nil)
	}()

	return closed
}

func NewGokitSrv(logger log.Logger) addsvcpb.AddServer {
	metricsObj := internal.NewMetrics()
	tracer := stdopentracing.GlobalTracer()

	// 依次创建 svc，endpoint，transport三层的对象，每一层都会在上一层基础上封装
	// 在svc和endpoint层以中间件的形式添加【指标上传、api日志】功能

	// service需要的所有对象都通过New传入
	svc := service.New(logger, metricsObj.Ints, metricsObj.Chars)
	// 在endpoint层和transport层添加路径追踪功能
	endpoints := endpoint.New(svc, logger, metricsObj.Duration, tracer)
	addSrv := transport.NewgRPCServer(endpoints, tracer, logger)
	return addSrv
}
