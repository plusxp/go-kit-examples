package grpc

import (
	"context"
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
	return &pb.SayHiResponse{
		BaseRsp: &pbcommon.BaseRsp{ErrCode: rsp.ErrCode},
		Reply:   rsp.Reply,
	}, nil
}

func (g *grpcServer) SayHi(ctx context1.Context, req *pb.SayHiRequest) (*pb.SayHiResponse, error) {
	_, rep, err := g.sayHi.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SayHiResponse), nil
}

func makeMakeADateHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.MakeADateEndpoint, decodeMakeADateRequest, encodeMakeADateResponse, options...)
}

func decodeMakeADateRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.MakeADateRequest)
	return &endpoint.MakeADateRequest{P1: req}, nil
}

func encodeMakeADateResponse(_ context.Context, r interface{}) (interface{}, error) {
	rsp := r.(*endpoint.MakeADateResponse)
	return rsp.P0, nil
}
func (g *grpcServer) MakeADate(ctx context1.Context, req *pb.MakeADateRequest) (*pb.MakeADateResponse, error) {
	_, rep, err := g.makeADate.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.MakeADateResponse), nil
}

func makeUpdateUserInfoHandler(endpoints endpoint.Endpoints, options []grpc.ServerOption) grpc.Handler {
	return grpc.NewServer(endpoints.UpdateUserInfoEndpoint, decodeUpdateUserInfoRequest, encodeUpdateUserInfoResponse, options...)
}

func decodeUpdateUserInfoRequest(_ context.Context, req interface{}) (interface{}, error) {
	r := req.(*pb.UpdateUserInfoRequest)
	return &endpoint.UpdateUserInfoRequest{P1: r}, nil
}

func encodeUpdateUserInfoResponse(_ context.Context, rsp interface{}) (interface{}, error) {
	r := rsp.(*endpoint.UpdateUserInfoResponse)
	return r.P0, nil
}
func (g *grpcServer) UpdateUserInfo(ctx context1.Context, req *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoResponse, error) {
	_, rep, err := g.updateUserInfo.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UpdateUserInfoResponse), nil
}
