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

	// new一个写好的RPC客户端, 这里演示了一种logger传递方法，实际项目中并不一定需要这样做
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
*/
func (gw *MyGateWay) MakeADate(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	c := helloclient.MustNew(gw.RawLogger())

	var err error
	var rsp *pb.MakeADateReply

	rsp, err = c.MakeADate(context.Background(), &pb.MakeADateRequest{
		BaseReq: &pbcommon.BaseReq{Plat: pbcommon.Plat_pc},
		DateStr: v["date"],
		WantSay: v["want_say"],
	})

	if err != nil {
		gw.Log("RPC MakeADate err", err)
		rsp = &pb.MakeADateReply{
			BaseRsp: &pbcommon.BaseRsp{ErrCode: pbcommon.R_RPC_ERR},
		}
	}
	gw.JSON(w, rsp)
}

/*
1. 第三个接口在网关中实现了身份验证以及统一从http body中反序列化rpc接口参数的实现
2. 这种方式在项目中可能应用的更多，它不是一种REST-FUL接口风格了，为了提高开发效率和可维护性，
	所有接口仅支持POST方式调用，token放在header中，请求参数放在body中传递，服务端仅从body中读取参数，
	然后直接透传至RPC接口
*/
func (gw *MyGateWay) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	// 不再从url中获取参数
	// v := mux.Vars(r)

	c := helloclient.MustNew(gw.RawLogger())

	rpcReq := &pb.UpdateUserInfoRequest{
		BaseReq: &pbcommon.BaseReq{},
	}
	// Prepare执行身份验证，反序列化等操作（从http header中读取token，从http body中读取rpc request args）
	ok := gw.Prepare(w, r, rpcReq)
	if !ok {
		// gateway.Prepare内部已经打了日志，这里不需要再log
		return
	}

	var err error
	var rsp *pb.UpdateUserInfoReply

	rsp, err = c.UpdateUserInfo(context.Background(), rpcReq)
	if err != nil {
		rsp = &pb.UpdateUserInfoReply{
			BaseRsp: &pbcommon.BaseRsp{ErrCode: pbcommon.R_RPC_ERR},
		}
		gw.Log("RPC UpdateUserInfo err", err)
	}
	gw.JSON(w, rsp)
}
