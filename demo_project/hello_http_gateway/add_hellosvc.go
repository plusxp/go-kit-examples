package main

import (
	"context"
	"github.com/gorilla/mux"
	grpcclient "hello/client/grpc"
	"hello/pb/gen-go/pb"
	"hello/pb/gen-go/pbcommon"
	"net/http"
)

/*
API defined here.
*/

// restful-api
func (gw *MyGateWay) SayHi(w http.ResponseWriter, r *http.Request) {

	// Extract params from path
	v := mux.Vars(r)

	// There is no need to verify params,
	name := v["name"]

	c := grpcclient.New()
	reply, code := c.SayHi(context.Background(), name)
	rsp := &pb.SayHiReply{
		BaseRsp: &pbcommon.BaseRsp{ErrCode: code},
		Reply:   reply,
	}
	gw.JSON(w, rsp)
}
