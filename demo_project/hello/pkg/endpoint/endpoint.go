package endpoint

import (
	"context"
	endpoint "github.com/go-kit/kit/endpoint"
	"hello/pb/gen-go/pbcommon"
	service "hello/pkg/service"
)

type SayHiRequest struct {
	Name string `json:"name"`
}

type SayHiResponse struct {
	Reply   string     `json:"reply"`
	ErrCode pbcommon.R `json:"err_code"`
}

func MakeSayHiEndpoint(s service.HelloService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*SayHiRequest)
		reply, err := s.SayHi(ctx, req.Name)
		return &SayHiResponse{
			ErrCode: err,
			Reply:   reply,
		}, nil
	}
}

func (e Endpoints) SayHi(ctx context.Context, name string) (reply string, errCode pbcommon.R) {
	request := &SayHiRequest{Name: name}
	response, err := e.SayHiEndpoint(ctx, request)
	if err != nil {
		return "", pbcommon.R_RPC_ERR
	}
	return response.(*SayHiResponse).Reply, response.(*SayHiResponse).ErrCode
}
