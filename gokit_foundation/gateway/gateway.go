package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	gogojsonpb "github.com/gogo/protobuf/jsonpb"
	gogoproto "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"gokit_foundation"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

const gatewayErrPrefix = "gateway:"

//// 网关域内的err类型
//type ErrCodeT int
//
//type ErrT struct {
//	Code ErrCodeT
//	Err  error
//}
//
///*
//网关域内的err_code常量，应该与于proto/common/resultcode.proto内的定义一致
//但不需要全部包括，根据需要定义
//*/
//const (
//	ECode_NoErr         ErrCodeT = iota
//	ECode_JsonUnmarshal ErrCodeT = 7
//	ECode_Auth          ErrCodeT = 8
//)

type Gateway struct {
	r *mux.Router
	*gokit_foundation.Logger
	addr string
}

func New(r *mux.Router, addr string, logger *gokit_foundation.Logger) *Gateway {
	return &Gateway{
		r:      r,
		Logger: logger,
		addr:   addr,
	}
}

func (g *Gateway) onStart() {
	g.Log("Gateway.OnStart:http-addr", g.addr)
	g.setupMW()
}

func (g *Gateway) onStop() {
	g.Log("Gateway.OnStop:http-addr", g.addr)
}

// Run使用最简洁的方式实现 统一start，优雅stop
func (g *Gateway) Run() error {
	g.onStart()
	defer g.onStop()

	srv := &http.Server{
		Addr: g.addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      g.r,
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	var err error
	go func() {
		err = http.ListenAndServe(g.addr, g.r)
	}()

	s := <-sc

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_ = srv.Shutdown(ctx)

	g.Log("Gateway.Run", "Stopped", "Signal", s)
	return err
}

// setupMW 安装中间件
func (g *Gateway) setupMW() {
	recoverMW := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					g.Log("Gateway.recover_mw", "==================== PANIC ====================")
					err = g.Log("Gateway.recover_mw", err)
					err = g.Log("Gateway.recover_mw", "===============================================")
					// 有些时候仍然有必要打印堆栈信息
					debug.PrintStack()
				}
			}()
			next.ServeHTTP(w, req)
		})
	}

	/*
		mux 使用mw的顺序与安装的顺序相反
	*/

	g.r.Use(recoverMW)
}

var defProtojsonpbM = &jsonpb.Marshaler{
	EnumsAsInts:  true, // proto enum类型仍然转为整数，而不是代表其含义的str
	EmitDefaults: true, // 仍然展示零值字段
	OrigName:     true,
}

var defGogojsonpbM = &gogojsonpb.Marshaler{
	EnumsAsInts:  true,
	EmitDefaults: true,
	OrigName:     true,
}

// JSON 直接封装+响应json数据.
func (g *Gateway) JSON(w http.ResponseWriter, rsp interface{}) error {
	var (
		b   []byte
		err error
	)
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			g.Log("Gateway.JSON: err", err)
			return
		}
		_, _ = w.Write(b)
	}()

	// 使用对应proto库的jsonpb工具来封装json，解决零值字段被忽略的问题
	// 一般来说一个项目要么使用标准protobuf库，要么使用gogo protobuf库
	// 了解gogo protobuf：https://github.com/gogo/protobuf

	buf := bytes.NewBuffer(nil)
	switch rsp.(type) {
	case proto.Message:
		err = defProtojsonpbM.Marshal(buf, rsp.(proto.Message))
		b = buf.Bytes()
	case gogoproto.Message:
		err = defGogojsonpbM.Marshal(buf, rsp.(gogoproto.Message))
		b = buf.Bytes()
	default:
		b, err = json.Marshal(rsp)
	}

	return err
}

// 用于传递logger
func (g *Gateway) RawLogger() *gokit_foundation.Logger {
	return g.Logger
}

// Prepare执行身份验证，反序列化等操作
// 当需要返回非200http状态码时，rsp.body为空即可
func (g *Gateway) Prepare(w http.ResponseWriter, httpReq *http.Request, rpcReq interface{}) (ok bool) {
	var err error
	defer func() {
		if err == nil {
			ok = true
			return
		}
		g.Log("Gateway.Prepare: err", err, "method", httpReq.Method, "path", httpReq.URL.Path)
	}()

	err = g.authenticate(httpReq)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if rpcReq != nil {
		if err = g.unmarshalReq(httpReq, rpcReq); err != nil {
			g.Log("Gateway.Prepare unmarshalReq err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	return
}

// 反序列化req
func (g *Gateway) unmarshalReq(httpReq *http.Request, req interface{}) error {
	return json.NewDecoder(httpReq.Body).Decode(req)
}

// 对req做身份验证
func (g *Gateway) authenticate(httpReq *http.Request) error {
	// JWT示例
	return authViaJwt(httpReq)
}

// 网关域内的err应该有特定前缀
func (g *Gateway) wrapErr(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s %v", gatewayErrPrefix, err)
}
