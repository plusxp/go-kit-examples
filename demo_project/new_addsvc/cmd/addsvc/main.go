package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"gokit_foundation"
	"os"
	"os/signal"
	"syscall"
)

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
	logger log.Logger
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

	beforeStart()

	// mainCtx传给每个后台运行的goroutine，mainCtx结束时，所有接收mainCtx的goroutine都应该退出
	mainCtx, stopBackgroundTasks := context.WithCancel(context.TODO())
	httpSrvClosed := startHttpSrv(mainCtx, httpSrvAddr)
	gRPCSrvClosed := startgRPCSrv(mainCtx, grpcSrvAddr)

	afterStart(srvHost, *grpcPort)

	logger.Log("gokit", "*******************Everything is ready*******************")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Log("............recv-signal", sig)

	{
		beforeStop()
		stopBackgroundTasks() // 停止所有后台任务
		afterStop()
	}
	{ // 等待各个服务正常退出
		<-httpSrvClosed
		<-gRPCSrvClosed
	}
	logger.Log("gokit", "exited!")
}
