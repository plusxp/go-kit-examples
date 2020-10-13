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

func newHelloClient(logger *gokit_foundation.Logger) *Client {
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
	// 出错时直接在这一层panic，外面就不需要处理
	logger.Must(err == nil, "HelloClient is nil")

	sc, err = NewSvc(conn)
	_util.PanicIfErr(err, nil)

	return &Client{
		HelloService: sc,
		conn:         conn,
	}
}

func MustNew(logger *gokit_foundation.Logger) *Client {
	if svcClient == nil || svcClient.conn == nil {
		svcClient = newHelloClient(logger)
	}
	return svcClient
}

func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
		// 允许Close后再MustNew
		c.conn = nil
	}
}
