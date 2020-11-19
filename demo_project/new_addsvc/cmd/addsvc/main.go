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
	"net"
	"net/http"
	"new_addsvc/internal"
	"new_addsvc/pb/gen-go/addsvcpb"
	"new_addsvc/pkg/endpoint"
	"new_addsvc/pkg/service"
	"new_addsvc/pkg/transport"
	"os"
	"time"
)

func NewAddSrv(logger log.Logger) addsvcpb.AddServer {
	metricsObj := internal.NewMetrics()
	tracer := stdopentracing.GlobalTracer()

	// 依次创建 svc，endpoint，transport三层的对象，每一层都会在上一层基础上封装
	// 在svc和endpoint层以中间件的形式添加【指标上传、api日志】功能

	// service需要的所有对象都通过New传入
	svc := service.New(logger, metricsObj.Ints, metricsObj.Chars)
	// 在endpoint层和transport层添加路径追踪功能
	endpoints := endpoint.New(svc, logger, metricsObj.Duration, tracer)
	addSrv := transport.NewGRPCServer(endpoints, tracer, logger)
	return addSrv
}

// for test
func init() {
	// 配置consul服务地址，必须是一个有效的consul地址
	os.Setenv("CONSUL_ADDR", "localhost:8500")
}

/*
new_addsvc服务依赖了一些外部中间件如下：
-	强依赖(若连不上则无法启动)
	-	consul
	-	redis
-	弱依赖(不需要连接或连不上也能启动)
	-	prometheus
*/

var (
	grpcSrv *grpc.Server
	httpSrv *http.Server
	logger  log.Logger
)

func main() {
	// 服务运行的主机地址，必须能够被你的consul-server访问，否则consul的健康检查会失败
	srvHost := "127.0.0.1"
	var grpcPort = flag.Int("grpc.port", 8080, "grpc listen address")
	var httpPort = flag.Int("http.port", 8081, "http listen address")

	grpcSrvAddr := fmt.Sprintf("%s:%d", srvHost, *grpcPort)
	httpSrvAddr := fmt.Sprintf("%s:%d", srvHost, *httpPort)

	flag.Parse()
	logger = gokit_foundation.NewKvLogger(nil)

	grpcSrv = grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	httpSrv = &http.Server{}

	/*
		这里使用 TaskGroup 完成程序的多任务同时启动，同时退出
		在实际项目中可以参考其思路，自行实现
	*/

	// 初始化一个TaskGroup对象
	tg := _go.NewTaskGroup()

	addTaskListenSignal(tg)
	initFirstly(srvHost, *grpcPort)

	addTaskHttpSrv(tg, httpSrvAddr)
	addTaskGRPCSrv(tg, grpcSrvAddr)

	logger.Log("main", "started")
	tg.Run()
}

// 添加后台任务：监听退出信号（第一个添加）
func addTaskListenSignal(tg *_go.TaskGroup) {
	tk, sc := _util.ListenSignalTask(logger, onClose)
	tg.Add(tk).Interrupt(func(err error) {
		close(sc)
	})
}

func addTaskHttpSrv(tg *_go.TaskGroup, httpSrvAddr string) {
	// http服务监听8080, 目前只提供metric接口给prometheus调用
	httpSrvTask := func(_ context.Context) error {
		logger.Log("NewTaskGroup", "httpSrvTask", "httpSrvAddr", httpSrvAddr)

		httpLis, err := net.Listen("tcp", httpSrvAddr)
		_util.PanicIfErr(err, nil)

		// default use http.DefaultServeMux as handler
		err = httpSrv.Serve(httpLis)
		return err
	}
	tg.Add(httpSrvTask).Interrupt(func(err error) {
		if err != nil {
			logger.Log("httpSrvTask", "exited", "err", err)
		} else {
			closeCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
			err := httpSrv.Shutdown(closeCtx)
			logger.Log("httpSrvTask", "exited", "clean", err)
		}
	})
}

func addTaskGRPCSrv(tg *_go.TaskGroup, grpcSrvAddr string) {
	// 添加后台任务：启动rpc-srv
	grpcSrvTask := func(_ context.Context) error {
		logger.Log("NewTaskGroup", "grpcSrvTask", "grpcSrvAddr", grpcSrvAddr)

		grpcLis, err := net.Listen("tcp", grpcSrvAddr)
		_util.PanicIfErr(err, nil)

		addSrv := NewAddSrv(logger)
		addsvcpb.RegisterAddServer(grpcSrv, addSrv)
		// 这里注册了AddSrv以及healthSrv
		gokit_foundation.RegisterGRPCHealthSrv(grpcSrv)

		err = grpcSrv.Serve(grpcLis)
		return err
	}
	tg.Add(grpcSrvTask).Interrupt(func(err error) {
		if err != nil {
			logger.Log("grpcSrvTask", "exited", "err", err)
		} else {
			grpcSrv.GracefulStop()
			logger.Log("grpcSrvTask", "exited", "clean", nil)
		}
	})
}
