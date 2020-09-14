package endpoint

import (
	"context"
	service "hello/pkg/service"

	endpoint "github.com/go-kit/kit/endpoint"
)

type Failure interface {
	Failed() error
}

type SayHiRequest struct {
	Name string `json:"name"`
}

type SayHiResponse struct {
	Reply string `json:"reply"`
	Err   error  `json:"err"`
}

func (r SayHiResponse) Failed() error {
	return r.Err
}

func MakeSayHiEndpoint(s service.HelloService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*SayHiRequest)
		reply, err := s.SayHi(ctx, req.Name)
		return &SayHiResponse{
			Err:   err,
			Reply: reply,
		}, nil
	}
}

func (e Endpoints) SayHi(ctx context.Context, name string) (reply string, err error) {
	request := &SayHiRequest{Name: name}
	response, err := e.SayHiEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(*SayHiResponse).Reply, response.(*SayHiResponse).Err
}
