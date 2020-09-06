package endpoint

import (
	"context"
	service "hello/pkg/service"

	endpoint "github.com/go-kit/kit/endpoint"
)

// SayHiRequest collects the request parameters for the SayHi method.
type SayHiRequest struct {
	Name string `json:"name"`
	Say  string `json:"say"`
}

// SayHiResponse collects the response parameters for the SayHi method.
type SayHiResponse struct {
	Reply string `json:"reply"`
	Err   error  `json:"err"`
}

// MakeSayHiEndpoint returns an endpoint that invokes SayHi on the service.
func MakeSayHiEndpoint(s service.HelloService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SayHiRequest)
		reply, err := s.SayHi(ctx, req.Name, req.Say)
		return SayHiResponse{
			Err:   err,
			Reply: reply,
		}, nil
	}
}

// Failed implements Failer.
func (r SayHiResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// SayHi implements Service. Primarily useful in a client.
func (e Endpoints) SayHi(ctx context.Context, name string, say string) (reply string, err error) {
	request := SayHiRequest{
		Name: name,
		Say:  say,
	}
	response, err := e.SayHiEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(SayHiResponse).Reply, response.(SayHiResponse).Err
}
