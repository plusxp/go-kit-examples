package endpoint

import (
	"context"
	endpoint "github.com/go-kit/kit/endpoint"
	"hello/pb/gen-go/pb"
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
	// 这个err不是svc返回的，而是封装了多个mw的endpoint返回的，属于意料之外的err，此时response可能是nil
	if err != nil {
		return "", pbcommon.R_RPC_ERR
	}
	return response.(*SayHiResponse).Reply, response.(*SayHiResponse).ErrCode
}

// MakeADateRequest collects the request parameters for the MakeADate method.
type MakeADateRequest struct {
	P1 *pb.MakeADateRequest `json:"p1"`
}

// MakeADateResponse collects the response parameters for the MakeADate method.
type MakeADateResponse struct {
	P0  *pb.MakeADateResponse `json:"p0"`
	Err error
}

// MakeMakeADateEndpoint returns an endpoint that invokes MakeADate on the service.
func MakeMakeADateEndpoint(s service.HelloService) endpoint.Endpoint {
	return func(c0 context.Context, request interface{}) (interface{}, error) {
		req := request.(*MakeADateRequest)
		p0, err := s.MakeADate(c0, req.P1)
		return &MakeADateResponse{P0: p0, Err: err}, err
	}
}

// MakeADate implements Service. Primarily useful in a client.
func (e Endpoints) MakeADate(c0 context.Context, p1 *pb.MakeADateRequest) (p0 *pb.MakeADateResponse, err error) {
	request := &MakeADateRequest{P1: p1}
	response, err := e.MakeADateEndpoint(c0, request)
	if err != nil {
		// 注意：endpoint层返回的err会被client端使用的限流、断路器等设施捕获
		// 所以这个err只有在系统级（如db故障，依赖服务异常）异常时返回，业务err应该放在reply中
		return nil, err
	}
	return response.(*MakeADateResponse).P0, err
}

// Failed implements Failer.
func (r MakeADateResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}

// UpdateUserInfoRequest collects the request parameters for the UpdateUserInfo method.
type UpdateUserInfoRequest struct {
	P1 *pb.UpdateUserInfoRequest `json:"p1"`
}

// UpdateUserInfoResponse collects the response parameters for the UpdateUserInfo method.
type UpdateUserInfoResponse struct {
	P0 *pb.UpdateUserInfoResponse `json:"p0"`
	E1 error                      `json:"e1"`
}

// MakeUpdateUserInfoEndpoint returns an endpoint that invokes UpdateUserInfo on the service.
func MakeUpdateUserInfoEndpoint(s service.HelloService) endpoint.Endpoint {
	return func(c0 context.Context, request interface{}) (interface{}, error) {
		req := request.(*UpdateUserInfoRequest)
		p0, e1 := s.UpdateUserInfo(c0, req.P1)
		return &UpdateUserInfoResponse{
			E1: e1,
			P0: p0,
		}, nil
	}
}

// Failed implements Failer.
func (r UpdateUserInfoResponse) Failed() error {
	return r.E1
}

// UpdateUserInfo implements Service. Primarily useful in a client.
func (e Endpoints) UpdateUserInfo(c0 context.Context, p1 *pb.UpdateUserInfoRequest) (p0 *pb.UpdateUserInfoResponse, e1 error) {
	request := &UpdateUserInfoRequest{P1: p1}
	response, err := e.UpdateUserInfoEndpoint(c0, request)
	if err != nil {
		return nil, err
	}
	return response.(*UpdateUserInfoResponse).P0, response.(*UpdateUserInfoResponse).E1
}
