package addtransport

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"go-kit-examples/new-to-gokit/5.addsvc_new/pb"
	"go-kit-examples/new-to-gokit/5.addsvc_new/pkg/addendpoint"
	"go-kit-examples/new-to-gokit/5.addsvc_new/pkg/addservice"
)

type grpcServer struct {
	sum    grpctransport.Handler
	concat grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer.
func NewGRPCServer(epMap addendpoint.EpMap, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer, logger log.Logger) pb.AddServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	if zipkinTracer != nil {
		// Zipkin GRPC Server Trace can either be instantiated per gRPC method with a
		// provided operation name or a global tracing service can be instantiated
		// without an operation name and fed to each Go kit gRPC server as a
		// ServerOption.
		// In the latter case, the operation name will be the endpoint's grpc method
		// path if used in combination with the Go kit gRPC Interceptor.
		//
		// In this example, we demonstrate a global Zipkin tracing service with
		// Go kit gRPC Interceptor.
		options = append(options, zipkin.GRPCServerTrace(zipkinTracer))
	}

	var methods = []string{"Sum", "Concat"}
	var notImplemented []string
	for _, s := range methods {
		if _, ok := epMap[s]; !ok {
			notImplemented = append(notImplemented, s)
		}
	}

	if len(notImplemented) > 0 {
		panic(fmt.Sprintf("NewGRPCServer: method(s) not implemented ---> %v", notImplemented))
	}

	//checkTokenBefore := func(ctx context.Context, md metadata.MD) context.Context{
	//	return ctx
	//}

	return &grpcServer{
		sum: grpctransport.NewServer(
			epMap["Sum"],
			decodeGRPCSumRequest,
			encodeGRPCSumResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Sum", logger)))...,
		),
		concat: grpctransport.NewServer(
			epMap["Concat"],
			decodeGRPCConcatRequest,
			encodeGRPCConcatResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Concat", logger)))...,
		),
	}
}

func (s *grpcServer) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumReply, error) {
	_, rep, err := s.sum.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SumReply), nil
}

func (s *grpcServer) Concat(ctx context.Context, req *pb.ConcatRequest) (*pb.ConcatReply, error) {
	_, rep, err := s.concat.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ConcatReply), nil
}

// decodeGRPCSumRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC sum request to a user-domain sum request. Primarily useful in a server.
func decodeGRPCSumRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SumRequest)
	return &addservice.SumRequest{A: int(req.A), B: int(req.B)}, nil
}

// decodeGRPCConcatRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC concat request to a user-domain concat request. Primarily useful in a
// server.
func decodeGRPCConcatRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ConcatRequest)
	return &addservice.ConcatRequest{A: req.A, B: req.B}, nil
}

// encodeGRPCSumResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain sum response to a gRPC sum reply. Primarily useful in a server.
func encodeGRPCSumResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*addservice.SumResponse)
	return &pb.SumReply{V: int64(resp.V), Err: err2str(resp.Err)}, nil
}

// encodeGRPCConcatResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain concat response to a gRPC concat reply. Primarily useful in a
// server.
func encodeGRPCConcatResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*addservice.ConcatResponse)
	return &pb.ConcatReply{V: resp.V, Err: err2str(resp.Err)}, nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
