package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"go-util/_go"
	"go-util/_util"
	"gokit_foundation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"net/http"
	config2 "new_addsvc/config"
	"new_addsvc/internal"
	"new_addsvc/pb/gen-go/addsvcpb"
	"new_addsvc/pkg/endpoint"
	"new_addsvc/pkg/service"
	"new_addsvc/pkg/transport"
	"os"
	"time"
)

func NewSvr(logger log.Logger) addsvcpb.AddServer {
	metricsObj := internal.NewMetrics()
	tracer := stdopentracing.GlobalTracer()

	// 依次创建 svc，endpoint，transport三层的对象，每一层都会在上一层基础上封装
	// 在svc和endpoint层以中间件的形式添加【指标上传、api日志】功能
	svc := service.New(logger, metricsObj.Ints, metricsObj.Chars)
	// 在endpoint层和transport层添加路径追踪功能
	endpoints := endpoint.New(svc, logger, metricsObj.Duration, tracer)
	grpcServer := transport.NewGRPCServer(endpoints, tracer, logger)
	return grpcServer
}

// for test
func init() {
	// 配置consul服务地址，必须是一个有效的consul地址
	os.Setenv("CONSUL_ADDR", "192.168.1.168:8500")
}

/*
new_addsvc服务依赖了一些外部中间件如下：
-	强依赖(若连不上则无法启动)
	-	consul
-	若依赖(不需要连接或连不上也能启动)
	-	prometheus
*/

func main() {
	// 这个地址必须能够被你的consul-server访问，否则consul的健康检查会失败
	svrHost := "127.0.0.1"
	var grpcPort = flag.Int("grpc.port", 8080, "grpc listen address")
	var httpPort = flag.Int("http.port", 8081, "http listen address")

	grpcSvrAddr := fmt.Sprintf("%s:%d", svrHost, *grpcPort)
	httpSvrAddr := fmt.Sprintf("%s:%d", svrHost, *httpPort)

	flag.Parse()

	logger := gokit_foundation.NewKvLogger(nil)
	grpcSvr := NewSvr(logger)

	grpcLis, err := net.Listen("tcp", grpcSvrAddr)
	_util.PanicIfErr(err, nil)

	httpLis, err := net.Listen("tcp", httpSvrAddr)
	_util.PanicIfErr(err, nil)

	/*
		初始化grpcSvr和httpSvr
	*/
	baseSvr := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	httpSvr := &http.Server{}

	// 创建一个所有后台任务共享的ctx，当进程退出时，所有后台任务都应该监听到ctx.Done()，然后graceful exit
	var uniformCtx, cancel = context.WithCancel(context.Background())

	/*
		这里使用 NewSafeAsyncTask 完成程序的一系列启动任务
		在实际项目中可以参考其思路，自行实现
	*/

	// 初始化一个SafeAsyncTask对象
	ak := _go.NewSafeAsyncTask(uniformCtx, cancel)

	// 程序退出时的操作
	onClose := func() {
		// 注意，首先应该先从consul删除实例信息，再执行其他操作
		gokit_foundation.ConsulDeregister()

		// 创建一个用于执行资源释放的ctx，避免时间过长
		// 这里因为程序要退出了，所以不用再 defer cancel(), 其他时候最好执行 defer cancel() 释放其内部资源
		closeCtx, _ := context.WithTimeout(context.Background(), time.Second*2)

		baseSvr.GracefulStop()
		err = httpSvr.Shutdown(closeCtx)
		_util.PanicIfErr(err, nil)
	}

	// 添加后台任务：监听退出信号（第一个添加）
	onSignalTask := _util.ListenSignalTask(uniformCtx, cancel, logger, onClose)
	ak.AddTask(onSignalTask)

	// 添加后台任务：service discovery
	svcRegTask := internal.SvcRegisterTask(uniformCtx, logger, config2.ServiceName, svrHost, *grpcPort)
	ak.AddTask(svcRegTask)

	// http服务监听8080, 目前只提供metric接口给prometheus调用
	httpSvc := func(_ context.Context, setter _go.Setter) {
		logger.Log("NewSafeAsyncTask", "httpSvc", "httpSvrAddr", httpSvrAddr)
		// default use http.DefaultServeMux as handler
		err := httpSvr.Serve(httpLis)
		// 调用SetErr后，若err!=nil，会使得所有task立即退出
		logger.Log("NewSafeAsyncTask", "httpSvc", "err", err)
		setter.SetErr(err)
	}
	ak.AddTask(httpSvc)

	// 添加后台任务：启动rpc-svr
	startGrpcSvrTask := func(_ context.Context, setter _go.Setter) {
		logger.Log("NewSafeAsyncTask", "GrpcSvr", "grpcSvrAddr", grpcSvrAddr)
		// 这里注册了AddSvr以及healthSvr
		addsvcpb.RegisterAddServer(baseSvr, grpcSvr)

		s := gokit_foundation.NewHealthCheckSvr()
		grpc_health_v1.RegisterHealthServer(baseSvr, s)

		err := baseSvr.Serve(grpcLis)
		if err != nil {
			logger.Log("NewSafeAsyncTask", "GrpcSvr", "err", err)
			setter.SetErr(err)
		}
	}
	ak.AddTask(startGrpcSvrTask)

	logger.Log("main", "started")
	ak.Run()
}
