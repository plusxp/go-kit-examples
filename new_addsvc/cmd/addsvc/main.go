package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"go-kit-examples/go-util/_go"
	"go-kit-examples/go-util/_util"
	"go-kit-examples/gokit_foundation"
	"go-kit-examples/new_addsvc/config"
	"go-kit-examples/new_addsvc/internal"
	"go-kit-examples/new_addsvc/pb/gen-go/addsvcpb"
	"go-kit-examples/new_addsvc/pkg/endpoint"
	"go-kit-examples/new_addsvc/pkg/service"
	"go-kit-examples/new_addsvc/pkg/transport"
	"google.golang.org/grpc"
	"net"
	"os"
)

func NewSvr(logger log.Logger) addsvcpb.AddServer {
	metricsObj := internal.NewMetrics()

	var tracer stdopentracing.Tracer
	tracer = stdopentracing.GlobalTracer()

	// grpctransport(endpoint(svc))
	svc := service.New(logger, metricsObj.Ints, metricsObj.Chars)
	endpoints := endpoint.New(svc, logger, metricsObj.Duration, tracer)
	grpcServer := transport.NewGRPCServer(endpoints, tracer, logger)
	return grpcServer
}

// for test
func init() {
	os.Setenv("CONSUL_ADDR", "192.168.1.168:8500")
}

func main() {
	// 这个地址必须能够被你的consul-server访问，否则consul的健康检查会失败
	svrHost := "192.168.1.10"
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

	svcRegTask := internal.SvcRegisterTask(ctx, config.ServiceName, svrHost, *port, baseServer)
	// 添加后台任务：service discovery
	ak.AddTask(svcRegTask...)

	startGrpcSvrTask := func(_ context.Context, setter _go.Setter) {
		logger.Log("transport", "gRPC", "addr", addr)
		addsvcpb.RegisterAddServer(baseServer, grpcServer)
		err := baseServer.Serve(grpcListener)
		setter.SetErr(err)
		// 调用SetErr后，若err!=nil，会使得所有task立即退出
	}

	// 添加后台任务：启动 rpc svr
	ak.AddTask(startGrpcSvrTask)

	ak.RunAndWait()
}
