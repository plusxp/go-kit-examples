package grpc

import (
	"context"
	"errors"
	"hello/pb/gen-go/pb"
	"hello/pb/gen-go/pbcommon"
	endpoint "hello/pkg/endpoint"

	grpc "github.com/go-kit/kit/transport/grpc"
	context1 "golang.org/x/net/context"
)

func makeSayHiHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.SayHiEndpoint, decodeSayHiRequest, encodeSayHiResponse, options...)
}

func decodeSayHiRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.SayHiRequest)
	return &endpoint.SayHiRequest{Name: req.Name}, nil
}

func encodeSayHiResponse(_ context.Context, r interface{}) (interface{}, error) {
	rsp := r.(*endpoint.SayHiResponse)
	return &pb.SayHiReply{
		BaseRsp: &pbcommon.BaseRsp{ErrCode: rsp.ErrCode},
		Reply:   rsp.Reply,
	}, nil
}

func (g *grpcServer) SayHi(ctx context1.Context, req *pb.SayHiRequest) (*pb.SayHiReply, error) {
	_, rep, err := g.sayHi.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SayHiReply), nil
}

func makeMakeADateHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.MakeADateEndpoint, decodeMakeADateRequest, encodeMakeADateResponse, options...)
}

func decodeMakeADateRequest(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Hello' Decoder is not impelemented")
}

func encodeMakeADateResponse(_ context.Context, r interface{}) (interface{}, error) {
	return nil, errors.New("'Hello' Encoder is not impelemented")
}
func (g *grpcServer) MakeADate(ctx context1.Context, req *pb.MakeADateRequest) (*pb.MakeADateReply, error) {
	_, rep, err := g.makeADate.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.MakeADateReply), nil
}
