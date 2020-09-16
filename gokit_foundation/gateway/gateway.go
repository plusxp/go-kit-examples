package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"go-util/_util"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

// Defines ExposedGateway is safe way to expose root gateway.
//type ExposedGateway interface {
//	Run() error
//	JSON(w http.ResponseWriter, rsp interface{}) error
//	log.Logger
//}

type Gateway struct {
	r      *mux.Router
	logger log.Logger
	addr   string
}

func New(r *mux.Router, addr string, logger log.Logger) *Gateway {
	return &Gateway{
		r:      r,
		logger: logger,
		addr:   addr,
	}
}

func (g *Gateway) onStart() {
	g.logger.Log("Gateway.OnStart:http-addr", g.addr)
	g.setupMW()
}

func (g *Gateway) onStop() {
	g.logger.Log("Gateway.OnStop:http-addr", g.addr)
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

	g.logger.Log("Gateway.Run", "Stopped", "Signal", s)
	return err
}

// setupMW setups mux middleware, you could customize that, it should be a private method.
func (g *Gateway) setupMW() {
	recoverMW := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					err := g.logger.Log("Gateway.recoverMW: got err", err)
					debug.PrintStack()
					_util.PanicIfErr(err, nil)
				}
			}()
			next.ServeHTTP(w, req)
		})
	}

	/*
		For mux pkg, the usage of middleware is used in the reverse order of installation
	*/

	g.r.Use(recoverMW)
}

// JSON is helper func to write response.
func (g *Gateway) JSON(w http.ResponseWriter, rsp interface{}) error {
	var (
		b   []byte
		err error
	)
	defer func() {
		if err != nil {
			g.logger.Log("Gateway.JSON: got err", err)
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

// Defines Err vars here is an example, you should define them in proto file if you use protobuf as data transport protocol.
var ErrReqParams = errors.New("ErrReqParams")
