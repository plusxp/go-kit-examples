package main

import (
	"context"
	"github.com/gorilla/mux"
	helloclient "hello/client/grpc"
	"hello/pb/gen-go/pb"
	"hello/pb/gen-go/pbcommon"
	"net/http"
)

/*
	代理hello服务
*/

/*
1. 第一种接口代理：完全还原service层的方法的入参出参(这就需要在网关层做参数的验证和封装)
2. 这里网关实现并没有过度封装handlerFunc，仍然把 w,r两个对象暴露给接口使用，把更多的自由留给开发者
*/
func (gw *MyGateWay) SayHi(w http.ResponseWriter, r *http.Request) {

	// 直接从url路径中提取参数
	v := mux.Vars(r)
	// 这个便利性来自于mux
	name := v["name"]

	// new一个写好的RPC客户端
	c := helloclient.MustNew(gw.RawLogger())

	// 像本地调用一样的远程调用
	reply, code := c.SayHi(context.Background(), name)

	rsp := &pb.SayHiReply{
		BaseRsp: &pbcommon.BaseRsp{ErrCode: code},
		Reply:   reply,
	}
	// JSON响应
	gw.JSON(w, rsp)
}

/*
1. 第二种接口定义方式也许更方便，在service，endpoint层直接使用pb协议定义好的req&rsp
2. 可完全将入参验证下放到service层
*/
func (gw *MyGateWay) MakeADate(w http.ResponseWriter, r *http.Request) {

	v := mux.Vars(r)
	dateStr := v["date"]

	// 先声明rsp
	rsp := &pb.MakeADateReply{
		BaseRsp: &pbcommon.BaseRsp{ErrCode: pbcommon.R_RPC_ERR},
	}

	defer func() {
		gw.JSON(w, rsp)
	}()

	c := helloclient.MustNew(gw.RawLogger())

	rpcRsp, err := c.MakeADate(context.Background(), &pb.MakeADateRequest{
		BaseReq: &pbcommon.BaseReq{Plat: pbcommon.Plat_pc},
		DateStr: dateStr,
		WantSay: "Do you willing to date with me?",
	})

	if err != nil {
		gw.Log("err", err)
		return
	}
	rsp = rpcRsp
}
