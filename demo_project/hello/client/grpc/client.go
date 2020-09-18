package grpc

import (
	"context"
	"go-util/_util"
	"google.golang.org/grpc"
	"hello/pkg/service"
	"log"
	"time"
)

type Client struct {
	service.HelloService
	conn *grpc.ClientConn
}

var svcClient *Client

func newSvcClient() *Client {
	var grpcOpts = []grpc.DialOption{
		grpc.WithInsecure(), // 因为没有使用tls，必须加上这个，否则连接失败
	}
	var err error
	var conn *grpc.ClientConn
	var sc service.HelloService

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	log.Println(111)
	conn, err = grpc.DialContext(ctx, "localhost:8082", grpcOpts...)
	_util.PanicIfErr(err, nil)
	log.Println(222)

	sc, err = NewSvc(conn)
	_util.PanicIfErr(err, nil)

	return &Client{
		HelloService: sc,
		conn:         conn,
	}
}

func New() *Client {
	if svcClient == nil {
		svcClient = newSvcClient()
	}
	return svcClient
}

func (c *Client) Stop() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
