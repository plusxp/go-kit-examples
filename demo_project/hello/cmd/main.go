package main

import service "hello/cmd/service"

// cd demo_project/
// kit g s hello -t grpc --dmw -p hello/pb/proto -i hello/pb/gen-go/pb

/*
hello服务依赖了一些外部中间件如下：
-	强依赖(若连不上则无法启动)
	-	consul
	-	redis
	-	mysql
-	弱依赖(不需要连接或连不上也能启动)
	-	prometheus
*/

func main() {
	service.Run()
}
