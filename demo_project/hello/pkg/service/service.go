package service

import "context"

// HelloService describes the service.
type HelloService interface {
	// Add your methods here
	SayHi(ctx context.Context, name string) (reply string, err error)
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

func (b *basicHelloService) SayHi(ctx context.Context, name string) (reply string, err error) {
	return "Hi," + name, err
}
