package grpc

import (
	"context"
	"errors"
	endpoint "hello/pkg/endpoint"
	pb "hello/pkg/grpc/pb"

	grpc "github.com/go-kit/kit/transport/grpc"
	context1 "golang.org/x/net/context"
)

// makeSayHiHandler creates the handler logic
func makeSayHiHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.SayHiEndpoint, decodeSayHiRequest, encodeSayHiResponse, options...)
}

// decodeSayHiResponse is a transport/grpc.DecodeRequestFunc that converts a
// gRPC request to a user-domain SayHi request.
// TODO implement the decoder
func decodeSayHiRequest(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Hello' Decoder is not impelemented")
}

// encodeSayHiResponse is a transport/grpc.EncodeResponseFunc that converts
// a user-domain response to a gRPC reply.
// TODO implement the encoder
func encodeSayHiResponse(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Hello' Encoder is not impelemented")
}
func (g *grpcServer) SayHi(ctx context1.Context, req *pb.SayHiRequest) (*pb.SayHiReply, error) {
	_, rep, err := g.sayHi.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SayHiReply), nil
}
