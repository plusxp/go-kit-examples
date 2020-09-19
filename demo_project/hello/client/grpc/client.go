package grpc

import (
	"context"
	"go-util/_util"
	"gokit_foundation"
	"google.golang.org/grpc"
	"hello/pkg/service"
	"time"
)

type Client struct {
	service.HelloService
	conn *grpc.ClientConn
}

var svcClient *Client

func newHelloClient(logger *gokit_foundation.UnionLogger) *Client {
	var grpcOpts = []grpc.DialOption{
		grpc.WithInsecure(), // 因为没有使用tls，必须加上这个，否则连接失败
		grpc.WithBlock(),    // 若不加这项，远程服务断开再恢复时，网关调用会继续失败
	}
	var err error
	var conn *grpc.ClientConn
	var sc service.HelloService

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conn, err = grpc.DialContext(ctx, "localhost:8082", grpcOpts...)
	if err != nil {
		_ = logger.Kvlgr.Log("client.go: newHelloClient.err", err)
		return nil
	}

	sc, err = NewSvc(conn)
	_util.PanicIfErr(err, nil)

	return &Client{
		HelloService: sc,
		conn:         conn,
	}
}

func New(logger *gokit_foundation.UnionLogger) *Client {
	if svcClient == nil {
		svcClient = newHelloClient(logger)
	}
	return svcClient
}

func (c *Client) Stop() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
