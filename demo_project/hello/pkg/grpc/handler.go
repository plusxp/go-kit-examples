package grpc

import (
	"context"
	pb "hello/pb"
	endpoint "hello/pkg/endpoint"

	grpc "github.com/go-kit/kit/transport/grpc"
	context1 "golang.org/x/net/context"
)

// makeSayHiHandler creates the handler logic
func makeSayHiHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.SayHiEndpoint, decodeSayHiRequest, encodeSayHiResponse, options...)
}

// decodeSayHiResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain SayHi request.
func decodeSayHiRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.SayHiRequest)
	return &endpoint.SayHiRequest{Name: req.Name}, nil
}

// encodeSayHiResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
func encodeSayHiResponse(_ context.Context, r interface{}) (interface{}, error) {
	rsp := r.(*endpoint.SayHiResponse)
	return &pb.SayHiReply{Reply: rsp.Reply}, nil
}

func (g *grpcServer) SayHi(ctx context1.Context, req *pb.SayHiRequest) (*pb.SayHiReply, error) {
	_, rep, err := g.sayHi.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SayHiReply), nil
}
