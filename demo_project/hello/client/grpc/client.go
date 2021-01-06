package grpc

import (
	"context"
	"fmt"
	"gokit_foundation"
	"google.golang.org/grpc"
	grpctransport "hello/pkg/grpc"
	"hello/pkg/service"
	"io"
	"time"
)

/*
没有使用服务注册发现的client
*/

type Client struct {
	service.HelloService
	conn        *grpc.ClientConn
	traceCloser io.Closer
}

var svcClient *Client

func MustNew(logger *gokit_foundation.Logger, upstreamSvc string) *Client {
	if svcClient == nil || svcClient.conn == nil {
		svcClient = newHelloClient(logger, upstreamSvc)
	}
	return svcClient
}

func newHelloClient(logger *gokit_foundation.Logger, upstreamSvc string) *Client {
	var grpcOpts = []grpc.DialOption{
		grpc.WithInsecure(), // 因为没有使用tls，必须加上这个，否则连接失败
		grpc.WithBlock(),    // 若不加这项，远程服务断开再恢复时，网关调用会继续失败
	}
	var err error
	var conn *grpc.ClientConn
	var sc service.HelloService
	var traceCloser io.Closer

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conn, err = grpc.DialContext(ctx, "localhost:8081", grpcOpts...)
	// 出错时直接在这一层panic，外面就不需要处理
	// logger.PanicIfErr 在panic的信息中添加了这一行代码的位置，在外层recover时会打印出来，以便快速定位panic行
	logger.PanicIfErr(err, nil, fmt.Sprintf("grpc.DialContext err:%v", err))

	sc, traceCloser, err = grpctransport.NewSvc(conn, upstreamSvc)
	logger.PanicIfErr(err, nil, fmt.Sprintf("grpctransport.NewSvc err:%v", err))

	return &Client{
		HelloService: sc,
		conn:         conn,
		traceCloser:  traceCloser,
	}
}

func (c *Client) Close() {
	c.traceCloser.Close()
	if c.conn != nil {
		_ = c.conn.Close()
		// 允许Close后再MustNew
		c.conn = nil
	}
}
