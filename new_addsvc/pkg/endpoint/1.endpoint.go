package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"go-kit-examples/new_addsvc/pb/gen-go/resultcode"
	"go-kit-examples/new_addsvc/pkg/service"
)

// eps 内包含的ep对应service每个接口，而且必须一致
// ep非常重要，在RPC调用时，client调用的也是ep，对调用者隐藏了transport层
type Endpoints struct {
	SumEndpoint    endpoint.Endpoint
	ConcatEndpoint endpoint.Endpoint
}

// 将一个Service对象转为Endpoints对象
func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		SumEndpoint:    MakeSumEndpoint(s),
		ConcatEndpoint: MakeConcatEndpoint(s),
	}
}

// 针对接口：Sum 的转换方法
func MakeSumEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SumRequest)
		v, err := s.Sum(ctx, req.A, req.B)
		// err 在logger mw打印出来 以debug
		// 谨慎处理endpoint层返回的err，这个err会被各种
		return SumResponse{RetCode: errToRetCode(err), V: v}, nil
	}
}

// 针对接口：Concat 的转换方法
func MakeConcatEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ConcatRequest)
		p, err := s.Concat(ctx, req.A, req.B)
		return ConcatResponse{RetCode: errToRetCode(err), V: p}, nil
	}
}

// 统一处理err
func errToRetCode(err error) resultcode.RESULT_CODE {
	switch err {
	case service.ErrIntOverflow, service.ErrMaxSizeExceeded, service.ErrTwoZeroes:
		return resultcode.RESULT_CODE_RET_INVALID_ARGS
	default:
		return resultcode.RESULT_CODE_RET_UNKNOWN_ERR
	}
}
