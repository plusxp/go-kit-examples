package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/circuitbreaker"
	stdendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	pb "hello/pb/gen-go/pb"
	endpoint2 "hello/pkg/endpoint"
	"time"
)

/*
GRPC-client侧
	每个endpoint都会按序安装跟踪、限速、断路器中间件
	go-kit的设计不是同一给所有接口安装，而是手动的给每一个接口安装，粒度细了，也多了一点代码量
*/

// 不是随便定的，是由proto文件中的package名和service名组合得到
// 在这里定义这个全局变量也许不是最终方案，待定~
var gRPCSvrName = "pb.Hello"

// client调用, 这个方法接收一个实例地址，以及中间件，然后创建出endpoint
func MakeClientEndpoints(instance string, otTracer stdopentracing.Tracer, logger log.Logger) (*endpoint2.Endpoints, error) {
	if instance == "" {
		return nil, errors.New("no instance")
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()

	conn, err := grpc.DialContext(ctx, instance, grpc.WithInsecure())
	if err != nil {
		// 这种情况很少发生，即从健康中心获得了一个健康的实例地址，却仍然连不上
		return nil, fmt.Errorf("failed: grpc.DialContext %s %s", instance, err)
	}

	return newGRPCClient(conn, otTracer, logger), nil
}

// newGRPCClient returns an HelloService backed by a gRPC server at the other end
// of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
// newGRPCClient 返回一个用grpc conn为底层连接的HelloService对象
func newGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, logger log.Logger) *endpoint2.Endpoints {
	var options []grpctransport.ClientOption

	// 每个endpoint安装统一的mw，也可以在wrappedEndpoint修改逻辑，根据method而设置不同的mw
	SayHiEndpoint := wrappedEndpoint(conn, otTracer, logger, options, "SayHi")
	MakeADateEndpoint := wrappedEndpoint(conn, otTracer, logger, options, "MakeADate")
	UpdateUserInfoEndpoint := wrappedEndpoint(conn, otTracer, logger, options, "UpdateUserInfo")

	return &endpoint2.Endpoints{
		SayHiEndpoint:          SayHiEndpoint,
		MakeADateEndpoint:      MakeADateEndpoint,
		UpdateUserInfoEndpoint: UpdateUserInfoEndpoint,
	}
}

func wrappedEndpoint(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, logger log.Logger, options []grpctransport.ClientOption, method string) stdendpoint.Endpoint {
	var ep stdendpoint.Endpoint

	enc, dec, grpcResp := getEncDecFuncByMethod(method)
	options = append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))

	ep = grpctransport.NewClient(conn, gRPCSvrName, method, enc, dec, grpcResp, options...).Endpoint()
	ep = opentracing.TraceClient(otTracer, "method")(ep)
	ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    method,
		Timeout: 10 * time.Second,
	}))(ep)

	return ep
}

func getEncDecFuncByMethod(method string) (grpctransport.EncodeRequestFunc, grpctransport.DecodeResponseFunc, interface{}) {
	var (
		enc      grpctransport.EncodeRequestFunc
		dec      grpctransport.DecodeResponseFunc
		grpcResp interface{}
	)
	switch method {
	case "SayHi":
		enc = encodeGRPCSayHiRequest
		dec = decodeGRPCSayHiResponse
		grpcResp = pb.SayHiResponse{}
	case "MakeADate":
		enc = encodeGRPCMakeADateRequest
		dec = decodeGRPCMakeADateResponse
		grpcResp = pb.MakeADateResponse{}
	case "UpdateUserInfo":
		enc = encodeGRPCUpdateUserInfoRequest
		dec = decodeGRPCUpdateUserInfoResponse
		grpcResp = pb.UpdateUserInfoResponse{}
	}
	return enc, dec, grpcResp
}

func encodeGRPCSayHiRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*endpoint2.SayHiRequest)
	return &pb.SayHiRequest{Name: req.Name}, nil
}

func decodeGRPCSayHiResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SayHiResponse)
	return &endpoint2.SayHiResponse{Reply: reply.Reply}, nil
}

func encodeGRPCMakeADateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*endpoint2.MakeADateRequest)
	return req.P1, nil
}

func decodeGRPCMakeADateResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.MakeADateResponse)
	return &endpoint2.MakeADateResponse{P0: reply}, nil
}

func encodeGRPCUpdateUserInfoRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*endpoint2.UpdateUserInfoRequest)
	return req.P1, nil
}

func decodeGRPCUpdateUserInfoResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UpdateUserInfoResponse)
	return &endpoint2.UpdateUserInfoResponse{P0: reply}, nil
}
