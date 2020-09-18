package grpc

import (
	"context"
	"github.com/go-kit/kit/circuitbreaker"
	endpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	grpc1 "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	grpc "google.golang.org/grpc"
	"hello/pb/gen-go/pb"
	endpoint1 "hello/pkg/endpoint"
	service "hello/pkg/service"
	"time"
)

// NewSvc returns an AddService backed by a gRPC server at the other end
//  of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewSvc(conn *grpc.ClientConn) (service.HelloService, error) {
	/*
		Create some security measures
	*/
	var otTracer stdopentracing.Tracer
	otTracer = stdopentracing.GlobalTracer()
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))
	breaker := circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "SayHi",
		Timeout: 30 * time.Second,
	}))

	// Create go-kit grpc hooks, e.g.
	//      - grpctransport.ClientAfter(),
	//      - grpctransport.ClientFinalizer()
	// Injecting tracing ctx to grpc metadata, optionally.
	grpcBefore := grpc1.ClientBefore(opentracing.ContextToGRPC(otTracer, log.NewNopLogger()))
	/*
		Install into endpoints with above measures
	*/
	var sayHiEndpoint endpoint.Endpoint
	{
		sayHiEndpoint = grpc1.NewClient(conn, "pb.Hello", "SayHi",
			encodeSayHiRequest, decodeSayHiResponse, pb.SayHiReply{}, grpcBefore).Endpoint()
		sayHiEndpoint = opentracing.TraceClient(otTracer, "sayHi")(sayHiEndpoint)
		sayHiEndpoint = limiter(sayHiEndpoint)
		sayHiEndpoint = breaker(sayHiEndpoint)
	}

	var makeADateEndpoint endpoint.Endpoint
	{
		makeADateEndpoint = grpc1.NewClient(conn, "pb.Hello", "MakeADate",
			encodeMakeADateRequest, decodeMakeADateResponse, pb.MakeADateReply{}, grpcBefore).Endpoint()
		makeADateEndpoint = opentracing.TraceClient(otTracer, "makeADate")(makeADateEndpoint)
		makeADateEndpoint = limiter(makeADateEndpoint)
		makeADateEndpoint = breaker(makeADateEndpoint)
	}

	return endpoint1.Endpoints{
		SayHiEndpoint:     sayHiEndpoint,
		MakeADateEndpoint: makeADateEndpoint,
	}, nil
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
	r := reply.(*pb.SayHiReply)
	return &endpoint1.SayHiResponse{Reply: r.Reply, ErrCode: r.BaseRsp.ErrCode}, nil
}

// encodeMakeADateRequest is a transport/grpc.EncodeRequestFunc that converts a
//  user-domain MakeADate request to a gRPC request.
func encodeMakeADateRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request.(*pb.MakeADateRequest), nil
}

// decodeMakeADateResponse is a transport/grpc.DecodeResponseFunc that converts
// a gRPC concat reply to a user-domain concat response.
func decodeMakeADateResponse(_ context.Context, reply interface{}) (interface{}, error) {
	return reply.(*pb.MakeADateReply), nil
}
