package main

import (
	"context"
	"github.com/gorilla/mux"
	helloclient "hello/client/grpc"
	"hello/pb/gen-go/pb"
	"hello/pb/gen-go/pbcommon"
	"net/http"
	"time"
)

/*
	代理hello服务
*/

// RESTFUL-API
// 这里网关实现并没有过度封装handlerFunc，仍然把 w,r两个对象暴露给接口使用
// 把更多的自由留给开发者
func (gw *MyGateWay) SayHi(w http.ResponseWriter, r *http.Request) {

	// 直接从url路径中提取参数
	v := mux.Vars(r)
	// 这个便利性来自于mux
	name := v["name"]

	// new一个写好的RPC客户端
	c := helloclient.New(gw.Logger())
	// 如果c==nil，直接panic即可，不要作隐藏处理（安装recover中间件让程序不要退出即可）
	// 也可以直接把Must放到 helloclient.New 里面去执行
	gw.Mylgr.Must(c != nil)

	// 像本地调用一样的远程调用
	reply, code := c.SayHi(context.Background(), name)

	rsp := &pb.SayHiReply{
		BaseRsp: &pbcommon.BaseRsp{ErrCode: code},
		Reply:   reply,
	}
	// JSON响应
	gw.JSON(w, rsp)
}

// 第二种接口定义方式也许更方便，在service，endpoint层直接使用pb协议定义好的req&rsp
// 作为方法的入参出参
func (gw *MyGateWay) MakeADate(w http.ResponseWriter, r *http.Request) {

	v := mux.Vars(r)
	date := v["date"]

	// 先声明rsp
	rsp := &pb.MakeADateReply{
		BaseRsp: &pbcommon.BaseRsp{ErrCode: pbcommon.R_RPC_ERR},
	}

	defer func() {
		gw.JSON(w, rsp)
	}()

	c := helloclient.New(gw.UnionLogger)
	// 如果c==nil，直接panic即可，不要作隐藏处理（安装recover中间件让程序不要退出即可）
	gw.Mylgr.Must(c != nil)

	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		rsp.BaseRsp.ErrCode = pbcommon.R_INVALID_ARGS
		return
	}

	rpcRsp, err := c.MakeADate(context.Background(), &pb.MakeADateRequest{
		BaseReq:  &pbcommon.BaseReq{Plat: pbcommon.Plat_pc},
		DateTime: t.Unix(),
		WantSay:  "Do you willing to date with me?",
	})

	if err != nil {
		gw.Kvlgr.Log("err", err)
		return
	}
	rsp = rpcRsp
}
