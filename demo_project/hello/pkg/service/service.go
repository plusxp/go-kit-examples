package service

import "context"

// HelloService describes the service.
type HelloService interface {
	// Add your methods here
	// e.x: Foo(ctx context.Context,s string)(rs string, err error)
	SayHi(ctx context.Context, name, say string) (reply string, err error)
}

type basicHelloService struct{}

func (b *basicHelloService) SayHi(ctx context.Context, name string, say string) (reply string, err error) {
	// TODO implement the business logic of SayHi
	return reply, err
}

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
