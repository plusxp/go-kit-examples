package endpoint

import (
	"context"
	"hello/pb/gen-go/pbcommon"
	service "hello/pkg/service"

	endpoint "github.com/go-kit/kit/endpoint"
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

// MakeADateRequest collects the request parameters for the MakeADate method.
type MakeADateRequest struct {
	P1 *pb.MakeADateReq `json:"p1"`
}

// MakeADateResponse collects the response parameters for the MakeADate method.
type MakeADateResponse struct {
	P0 *pb.MakeADateRsp `json:"p0"`
}

// MakeMakeADateEndpoint returns an endpoint that invokes MakeADate on the service.
func MakeMakeADateEndpoint(s service.HelloService) endpoint.Endpoint {
	return func(c0 context.Context, request interface{}) (interface{}, error) {
		req := request.(*MakeADateRequest)
		p0 := s.MakeADate(c0, req.P1)
		return &MakeADateResponse{P0: p0}, nil
	}
}

// MakeADate implements Service. Primarily useful in a client.
func (e Endpoints) MakeADate(c0 context.Context, p1 *pb.MakeADateReq) (p0 *pb.MakeADateRsp) {
	request := &MakeADateRequest{P1: p1}
	response, err := e.MakeADateEndpoint(c0, request)
	if err != nil {
		return
	}
	return response.(*MakeADateResponse).P0
}
