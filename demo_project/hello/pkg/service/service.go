package service

import (
	"context"
	"fmt"
	"hello/pb/gen-go/pb"
	"hello/pb/gen-go/pbcommon"
	"time"
)

// HelloService describes the service.
type HelloService interface {
	// 服务方法，对应接口
	SayHi(ctx context.Context, name string) (reply string, err pbcommon.R)

	// 为了方便client更直接的返回response，service层也可以直接使用定好的req&rsp协议
	MakeADate(context.Context, *pb.MakeADateReq) *pb.MakeADateRsp
}

type basicHelloService struct{}

// NewBasicHelloService returns a naive, stateless implementation of HelloService.
func NewBasicHelloService() HelloService {
	return &basicHelloService{}
}

// New returns a HelloService with all of the expected middleware wired in.
func New(middleware []Middleware) HelloService {
	var svc HelloService = NewBasicHelloService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}

func (b *basicHelloService) SayHi(ctx context.Context, name string) (reply string, err pbcommon.R) {
	if name == "XI" {
		return "", pbcommon.R_INVALID_ARGS
	}
	return "Hi," + name, err
}

func (b *basicHelloService) MakeADate(c0 context.Context, p1 *pb.MakeADateReq) (p0 *pb.MakeADateRsp) {
	t := time.Unix(p1.DateTime, 0)
	month, day := t.Month(), t.Day()

	p0 = &pb.MakeADateRsp{
		BaseRsp: &pbcommon.BaseRsp{},
		Reply:   fmt.Sprintf("Sorry, I am too busy~"),
	}

	// 只接受10月1日作为约会时间
	if month == 10 && day == 1 {
		p0.Reply = fmt.Sprintf("OK~, I was going to arrive on 10.1")
	}
	return p0
}
