package gateway

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-util/_util"
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
	*gokit_foundation.UnionLogger
	addr string
}

func New(r *mux.Router, addr string, logger *gokit_foundation.UnionLogger) *Gateway {
	return &Gateway{
		r:           r,
		UnionLogger: logger,
		addr:        addr,
	}
}

func (g *Gateway) onStart() {
	g.Kvlgr.Log("Gateway.OnStart:http-addr", g.addr)
	g.setupMW()
}

func (g *Gateway) onStop() {
	g.Kvlgr.Log("Gateway.OnStop:http-addr", g.addr)
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

	g.Kvlgr.Log("Gateway.Run", "Stopped", "Signal", s)
	return err
}

// setupMW 安装中间件
func (g *Gateway) setupMW() {
	recoverMW := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					g.Kvlgr.Log("Gateway.recover_mw", "==================== PANIC ====================")
					err = g.Kvlgr.Log("Gateway.recover_mw", err)
					err = g.Kvlgr.Log("Gateway.recover_mw", "===============================================")
					// 有些时候仍然有必要打印堆栈信息
					debug.PrintStack()
					_util.PanicIfErr(err, nil)
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

// JSON 直接封装+响应json数据.
func (g *Gateway) JSON(w http.ResponseWriter, rsp interface{}) error {
	var (
		b   []byte
		err error
	)
	defer func() {
		if err != nil {
			g.Kvlgr.Log("Gateway.JSON: got err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	b, err = json.Marshal(rsp)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// 用于传递logger
func (g *Gateway) Logger() *gokit_foundation.UnionLogger {
	return g.UnionLogger
}
