package grpc

import (
	"context"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	grpc1 "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	"gokit_foundation"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"hello/pb/gen-go/pb"
	endpoint1 "hello/pkg/endpoint"
	"hello/pkg/service"
	"io"
	"time"
)

// NewSvc returns an AddService backed by a gRPC server at the other end
//  of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewSvc(conn *grpc.ClientConn, svc string) (service.HelloService, io.Closer, error) {
	/*
		Create some security measures
	*/
	closer, err := gokit_foundation.InitTracer(svc)
	if err != nil {
		return nil, nil, err
	}

	otTracer := stdopentracing.GlobalTracer()

	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))
	breaker := func(method string) endpoint.Middleware {
		return circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    method,
			Timeout: 30 * time.Second,
		}))
	}

	var svcPBName = "pb.Hello"
	// Create go-kit grpc hooks, e.g.
	//      - grpctransport.ClientAfter(),
	//      - grpctransport.ClientFinalizer()
	// Injecting tracing ctx to grpc metadata, optionally.
	grpcBefore := grpc1.ClientBefore(opentracing.ContextToGRPC(otTracer, log.NewNopLogger()))
	/*
		Install into endpoints with above measures
	*/
	var SayHiEndpoint endpoint.Endpoint
	{
		SayHiEndpoint = grpc1.NewClient(conn, svcPBName, "SayHi",
			encodeSayHiRequest, decodeSayHiResponse, pb.SayHiResponse{}, grpcBefore).Endpoint()
		SayHiEndpoint = opentracing.TraceClient(otTracer, "SayHi")(SayHiEndpoint)
		SayHiEndpoint = limiter(SayHiEndpoint)
		SayHiEndpoint = breaker("SayHi")(SayHiEndpoint)
	}

	var MakeADateEndpoint endpoint.Endpoint
	{
		MakeADateEndpoint = grpc1.NewClient(conn, svcPBName, "MakeADate",
			encodeMakeADateRequest, decodeMakeADateResponse, pb.MakeADateResponse{}, grpcBefore).Endpoint()
		MakeADateEndpoint = opentracing.TraceClient(otTracer, "MakeADate")(MakeADateEndpoint)
		MakeADateEndpoint = limiter(MakeADateEndpoint)
		MakeADateEndpoint = breaker("MakeADate")(MakeADateEndpoint)
	}

	var UpdateUserInfoEndpoint endpoint.Endpoint
	{
		UpdateUserInfoEndpoint = grpc1.NewClient(conn, svcPBName, "UpdateUserInfo",
			encodeUpdateUserInfoRequest, decodeUpdateUserInfoResponse, pb.UpdateUserInfoResponse{}, grpcBefore).Endpoint()
		UpdateUserInfoEndpoint = opentracing.TraceClient(otTracer, "UpdateUserInfo")(UpdateUserInfoEndpoint)
		UpdateUserInfoEndpoint = limiter(UpdateUserInfoEndpoint)
		UpdateUserInfoEndpoint = breaker("UpdateUserInfo")(UpdateUserInfoEndpoint)
	}

	return endpoint1.Endpoints{
		SayHiEndpoint:          SayHiEndpoint,
		MakeADateEndpoint:      MakeADateEndpoint,
		UpdateUserInfoEndpoint: UpdateUserInfoEndpoint,
	}, closer, nil
}

// encodeSayHiRequest is a transport/grpc.EncodeRequestFunc that converts a
//  user-domain SayHi request to a gRPC request.
func encodeSayHiRequest(_ context.Context, request interface{}) (interface{}, error) {
	r := request.(*endpoint1.SayHiRequest)
	return &pb.SayHiRequest{Name: r.Name}, nil
}

// decodeSayHiResponse is a transport/grpc.DecodeResponseFunc that converts
// a gRPC concat reply to a user-domain concat response.
func decodeSayHiResponse(_ context.Context, reply interface{}) (interface{}, error) {
	r := reply.(*pb.SayHiResponse)
	return &endpoint1.SayHiResponse{Reply: r.Reply}, nil
}

// encodeMakeADateRequest is a transport/grpc.EncodeRequestFunc that converts a
//  user-domain MakeADate request to a gRPC request.
func encodeMakeADateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*endpoint1.MakeADateRequest)
	return req.P1, nil
}

// decodeMakeADateResponse is a transport/grpc.DecodeResponseFunc that converts
// a gRPC concat reply to a user-domain concat response.
func decodeMakeADateResponse(_ context.Context, reply interface{}) (interface{}, error) {
	rsp := reply.(*pb.MakeADateResponse)
	return &endpoint1.MakeADateResponse{P0: rsp}, nil
}

// encodeUpdateUserInfoRequest is a transport/grpc.EncodeRequestFunc that converts a
//  user-domain UpdateUserInfo request to a gRPC request.
func encodeUpdateUserInfoRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*endpoint1.UpdateUserInfoRequest)
	return req.P1, nil
}

// decodeUpdateUserInfoResponse is a transport/grpc.DecodeResponseFunc that converts
// a gRPC concat reply to a user-domain concat response.
func decodeUpdateUserInfoResponse(_ context.Context, reply interface{}) (interface{}, error) {
	rsp := reply.(*pb.UpdateUserInfoResponse)
	return &endpoint1.UpdateUserInfoResponse{P0: rsp}, nil
}
