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

// RESTFUL-API
// 这里网关的实现并没有过度封装handlerFunc，仍然把 w,r两个对象暴露给接口使用
// 把更多的自由留给开发者
func (gw *MyGateWay) SayHi(w http.ResponseWriter, r *http.Request) {

	// 直接从url路径中提取参数
	v := mux.Vars(r)
	// 这个便利性来自于mux
	name := v["name"]

	// new一个写好的RPC客户端
	c := grpcclient.New()

	// 像本地调用一样的远程调用
	reply, code := c.SayHi(context.Background(), name)

	rsp := &pb.SayHiReply{
		BaseRsp: &pbcommon.BaseRsp{ErrCode: code},
		Reply:   reply,
	}
	// JSON响应
	gw.JSON(w, rsp)
}
