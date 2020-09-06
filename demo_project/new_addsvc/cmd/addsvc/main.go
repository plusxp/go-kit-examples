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
	config2 "new_addsvc/config"
	internal2 "new_addsvc/internal"
	"new_addsvc/pb/gen-go/addsvcpb"
	endpoint2 "new_addsvc/pkg/endpoint"
	service2 "new_addsvc/pkg/service"
	transport2 "new_addsvc/pkg/transport"
	"os"
)

func NewSvr(logger log.Logger) addsvcpb.AddServer {
	// todo 启动http.DefaultServeMux
	metricsObj := internal2.NewMetrics()
	//http.ListenAndServe(":8080", nil)
	var tracer stdopentracing.Tracer
	tracer = stdopentracing.GlobalTracer()

	// grpctransport(endpoint(svc))
	svc := service2.New(logger, metricsObj.Ints, metricsObj.Chars)
	endpoints := endpoint2.New(svc, logger, metricsObj.Duration, tracer)
	grpcServer := transport2.NewGRPCServer(endpoints, tracer, logger)
	return grpcServer
}

// for test
func init() {
	os.Setenv("CONSUL_ADDR", "192.168.1.168:8500")
}

func main() {
	// 这个地址必须能够被你的consul-server访问，否则consul的健康检查会失败
	svrHost := "127.0.0.1"
	var port = flag.Int("grpc.port", 8080, "grpc listen address")

	addr := fmt.Sprintf("%s:%d", svrHost, *port)

	flag.Parse()

	logger := gokit_foundation.NewLogger()
	grpcServer := NewSvr(logger)

	grpcListener, err := net.Listen("tcp", addr)
	_util.PanicIfErr(err, nil)

	// 创建一个所有后台任务共享的ctx，当进程退出时，所有后台任务都应该监听到ctx.Done()，然后graceful exit
	var ctx, cancel = context.WithCancel(context.Background())

	// 初始化一个SafeAsyncTask对象
	ak := _go.NewSafeAsyncTask(ctx, cancel)

	// 程序退出时的操作
	onClose := func() {
		// 注意，首先应该先从consul删除实例信息，再执行其他操作
		gokit_foundation.ConsulDeregister()
		grpcListener.Close()
	}

	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))

	// 添加后台任务：监听退出信号（第一个添加）
	onSignalTask := _util.ListenSignalTask(ctx, cancel, logger, onClose)
	ak.AddTask(onSignalTask)

	svcRegTask := internal2.SvcRegisterTask(ctx, config2.ServiceName, svrHost, *port)
	// 添加后台任务：service discovery
	ak.AddTask(svcRegTask)

	startGrpcSvrTask := func(_ context.Context, setter _go.Setter) {
		logger.Log("transport", "gRPC", "addr", addr)
		// 这里注册了AddSvr以及healthSvr
		addsvcpb.RegisterAddServer(baseServer, grpcServer)

		s := gokit_foundation.NewHealthCheckSvr()
		grpc_health_v1.RegisterHealthServer(baseServer, s)

		err := baseServer.Serve(grpcListener)
		setter.SetErr(err)
		// 调用SetErr后，若err!=nil，会使得所有task立即退出
	}

	// 添加后台任务：启动 rpc svr
	ak.AddTask(startGrpcSvrTask)

	ak.RunAndWait()
}
