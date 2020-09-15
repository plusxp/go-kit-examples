package service

import (
	"context"
	"hello/pb/gen-go/pbcommon"
)

// HelloService describes the service.
type HelloService interface {
	// Add your methods here
	SayHi(ctx context.Context, name string) (reply string, err pbcommon.R)
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
