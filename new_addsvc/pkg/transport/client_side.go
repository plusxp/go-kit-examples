package transport

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/circuitbreaker"
	stdendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	pb "go-kit-examples/new_addsvc/pb/gen-go/addsvcpb"
	"go-kit-examples/new_addsvc/pkg/endpoint"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"os"
	"time"
)

/*
RPC-client侧
	每个接口都会按序安装跟踪、限速、断路器中间件
	go-kit的写法不是同一给所有接口安装，而是手动的给每一个接口安装，粒度细了，也多了一点代码量
*/

// client调用, 这个方法接收一个实例地址，以及中间件，然后创建出endpoint
func MakeClientEndpoints(instance string, otTracer stdopentracing.Tracer, logger log.Logger) (endpoint.AddSvcEndpoints, error) {
	if instance == "" {
		return endpoint.AddSvcEndpoints{}, errors.New("no instance")
	}

	// 这里可以设置dial选项
	conn, err := grpc.Dial(instance, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	return newGRPCClient(conn, otTracer, logger), nil
}

// newGRPCClient returns an AddService backed by a gRPC server at the other end
// of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func newGRPCClient(conn *grpc.ClientConn, otTracer stdopentracing.Tracer, logger log.Logger) endpoint.AddSvcEndpoints {
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	// global client middlewares
	var options []grpctransport.ClientOption

	// Each individual endpoint is an grpc/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	var sumEndpoint stdendpoint.Endpoint
	{
		sumEndpoint = grpctransport.NewClient(
			conn,
			"pb.Add",
			"Sum",
			encodeGRPCSumRequest,
			decodeGRPCSumResponse,
			pb.SumReply{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		sumEndpoint = opentracing.TraceClient(otTracer, "Sum")(sumEndpoint)
		sumEndpoint = limiter(sumEndpoint)
		sumEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Sum",
			Timeout: 30 * time.Second,
		}))(sumEndpoint)
	}

	// The Concat endpoint is the same thing, with slightly different
	// middlewares to demonstrate how to specialize per-endpoint.
	var concatEndpoint stdendpoint.Endpoint
	{
		concatEndpoint = grpctransport.NewClient(
			conn,
			"pb.Add",
			"Concat",
			encodeGRPCConcatRequest,
			decodeGRPCConcatResponse,
			pb.ConcatReply{},
			append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint()
		concatEndpoint = opentracing.TraceClient(otTracer, "Concat")(concatEndpoint)
		concatEndpoint = limiter(concatEndpoint)
		concatEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Concat",
			Timeout: 10 * time.Second,
		}))(concatEndpoint)
	}

	return endpoint.AddSvcEndpoints{
		SumEndpoint:    sumEndpoint,
		ConcatEndpoint: concatEndpoint,
	}
}

// decodeGRPCSumResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC sum reply to a user-domain sum response. Primarily useful in a client.
// 负责：endpointReq ==> grpcReq, client使用
func decodeGRPCSumResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SumReply)
	return endpoint.SumResponse{V: int(reply.V), RetCode: reply.Retcode}, nil
}

// decodeGRPCConcatResponse is a transport/grpc.DecodeResponseFunc that converts
// a gRPC concat reply to a user-domain concat response. Primarily useful in a
// client.
func decodeGRPCConcatResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ConcatReply)
	return endpoint.ConcatResponse{V: reply.V, RetCode: reply.Retcode}, nil
}

// encodeGRPCSumRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain sum request to a gRPC sum request. Primarily useful in a client.
// 负责：endpointReq ==> grpcReq, client使用
func encodeGRPCSumRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoint.SumRequest)
	return &pb.SumRequest{A: int64(req.A), B: int64(req.B)}, nil
}

// encodeGRPCConcatRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain concat request to a gRPC concat request. Primarily useful in a
// client.
func encodeGRPCConcatRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(endpoint.ConcatRequest)
	return &pb.ConcatRequest{A: req.A, B: req.B}, nil
}
