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
	MakeADate(context.Context, *pb.MakeADateRequest) (*pb.MakeADateReply, error)
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

// c0,p1是kit默认的变量命名规则，暂时认为没必要改
func (b *basicHelloService) MakeADate(c0 context.Context, p1 *pb.MakeADateRequest) (p0 *pb.MakeADateReply, err error) {
	p0 = &pb.MakeADateReply{
		BaseRsp: &pbcommon.BaseRsp{},
	}

	t, err := time.Parse("2006-01-02", p1.DateStr)
	if err != nil {
		p0.BaseRsp.ErrCode = pbcommon.R_INVALID_ARGS
		return
	}

	p0.Reply = fmt.Sprintf("Sorry, I am too busy~")

	month, day := t.Month(), t.Day()

	// 手动抛出错误，仍然应该正常返回rsp
	if month == 12 && day == 12 {
		return p0, fmt.Errorf("dependency svc err")
	}

	// 只接受10月1日作为约会时间
	if month == 10 && day == 1 {
		p0.Reply = fmt.Sprintf("OK~, I will arrive on 10.1")
	}
	return p0, nil
}
