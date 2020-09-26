package gateway

import (
	"bytes"
	"context"
	"encoding/json"
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
}

var defGogojsonpbM = &gogojsonpb.Marshaler{
	EnumsAsInts:  true,
	EmitDefaults: true,
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
			g.Log("Gateway.JSON: got err", err)
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
