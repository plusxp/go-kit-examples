package main

import service "hello/cmd/service"

// cd demo_project/
// kit g s hello -t grpc --dmw -p hello/pb/proto -i hello/pb/gen-go/pb

func main() {
	service.Run()
}
