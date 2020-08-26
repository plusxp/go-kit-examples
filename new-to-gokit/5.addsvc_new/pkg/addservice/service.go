package addservice

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"time"
)

type A func()

// Service describes a service that adds things together.
type Service interface {
	Sum(ctx context.Context, request interface{}) (response interface{}, err error)
	Concat(ctx context.Context, request interface{}) (response interface{}, err error)
}

// New returns a basic Service with all of the expected middlewares wired in.
func New(logger log.Logger) Service {
	var svc Service
	// 使用洋葱模式封装svc
	svc = NewBasicService(logger)
	return svc
}

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService(log log.Logger) Service {
	return &basicService{log}
}

type basicService struct {
	log.Logger
}

const (
	intMax = 1<<31 - 1
	intMin = -(intMax + 1)
	maxLen = 10
)

func (s basicService) Sum(ctx context.Context, req interface{}) (rsp interface{}, err error) {
	request := req.(*SumRequest)
	response := &SumResponse{}
	response.V = request.A + request.B
	time.Sleep(time.Millisecond * 100)
	//_ = s.Log("Sum", ctx.Value("k"))
	return response, nil
}

// Concat implements Service.
func (s basicService) Concat(_ context.Context, req interface{}) (rsp interface{}, err error) {
	request := req.(*ConcatRequest)
	response := &ConcatResponse{}
	response.V = request.A + request.B
	time.Sleep(time.Millisecond * 100)
	return response, nil
}

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = SumResponse{}
	_ endpoint.Failer = ConcatResponse{}
)

// SumRequest collects the request parameters for the Sum method.
type SumRequest struct {
	A, B int
}

// SumResponse collects the response values for the Sum method.
type SumResponse struct {
	V   int   `json:"v"`
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

// Failed implements endpoint.Failer.
func (r SumResponse) Failed() error { return r.Err }

// ConcatRequest collects the request parameters for the Concat method.
type ConcatRequest struct {
	A, B string
}

// ConcatResponse collects the response values for the Concat method.
type ConcatResponse struct {
	V   string `json:"v"`
	Err error  `json:"-"`
}

// Failed implements endpoint.Failer.
func (r ConcatResponse) Failed() error { return r.Err }
