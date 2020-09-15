package gateway

import (
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Defines ExposedGateway is safe way to expose root gateway.
type ExposedGateway interface {
	Run() error
	JSON(w http.ResponseWriter, rsp interface{}) error
	log.Logger
}

type Gateway struct {
	r      *mux.Router
	logger log.Logger
	addr   string
}

func New(r *mux.Router, addr string, logger log.Logger) ExposedGateway {
	return &Gateway{
		r:      r,
		logger: logger,
		addr:   addr,
	}
}

func (g *Gateway) OnStart() {
	g.logger.Log("Gateway.OnStart:http-addr", g.addr)
	g.setupMW()
}

func (g *Gateway) OnStop() {
	g.logger.Log("Gateway.OnStop:http-addr", g.addr)
}

func (g *Gateway) Run() error {
	defer g.OnStop()
	g.OnStart()

	sc := make(chan os.Signal)
	ch := make(chan error)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	go func() {
		ch <- http.ListenAndServe(g.addr, g.r)
	}()

	var err error
	select {
	case err = <-ch:
	case s := <-sc:
		g.logger.Log("Gateway.Run:on-signal", s)
	}
	return err
}

func (g *Gateway) Log(keyvals ...interface{}) error {
	return g.logger.Log(keyvals...)
}

// setupMW setups mux middleware, you could customize that, it should be a private method.
func (g *Gateway) setupMW() {
	recoverMW := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					_ = g.logger.Log("Gateway.recoverMW: got err", err)
				}
			}()
			next.ServeHTTP(w, req)
		})
	}
	g.r.Use(recoverMW)
	// Processing CORS issue is common.
	g.r.Use(mux.CORSMethodMiddleware(g.r))
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
		w.WriteHeader(http.StatusOK)
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
