package service

import (
	"flag"
	"fmt"
	"go-util/_str"
	"gokit_foundation"
	pb "hello/pb/gen-go/pb"
	endpoint "hello/pkg/endpoint"
	grpc "hello/pkg/grpc"
	service "hello/pkg/service"
	"net"
	http1 "net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	endpoint1 "github.com/go-kit/kit/endpoint"
	log "github.com/go-kit/kit/log"
	prometheus "github.com/go-kit/kit/metrics/prometheus"
	group "github.com/oklog/oklog/pkg/group"
	opentracinggo "github.com/opentracing/opentracing-go"
	prometheus1 "github.com/prometheus/client_golang/prometheus"
	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
	grpc1 "google.golang.org/grpc"
)

var tracer opentracinggo.Tracer
var logger *gokit_foundation.Logger

// Define our flags. Your service probably won't need to bind listeners for
// all* supported transports, but we do it here for demonstration purposes.
var fs = flag.NewFlagSet("hello", flag.ExitOnError)
var debugAddr = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
var grpcAddr = fs.String("grpc-addr", ":8081", "gRPC listen address")

func Run() {
	fs.Parse(os.Args[1:])

	s := strings.Split(*grpcAddr, ":")
	if len(s) != 2 {
		panic("Run: wrong grpcAddr")
	}
	grpcHost, grpcPort := s[0], s[1]
	if grpcHost == "" {
		grpcHost = "localhost"
	}
	// Create a single logger, which we'll use and give to other components.
	logger = gokit_foundation.NewLogger(nil)

	//  Determine which tracer to use. We'll pass the tracer to all the
	// components that use it, as a dependency

	logger.Log("tracer", "none")
	tracer = opentracinggo.GlobalTracer()

	// Init firstly before create service
	initFirstly(logger, grpcHost, _str.MustToInt(grpcPort))
	defer onClose()

	svc := service.New(getServiceMiddleware(logger), logger)
	eps := endpoint.New(svc, getEndpointMiddleware(logger))
	g := createService(eps)
	initMetricsEndpoint(g)
	initCancelInterrupt(g)

	logger.Log("exit", g.Run())
}

func initGRPCHandler(endpoints endpoint.Endpoints, g *group.Group) {
	options := defaultGRPCOptions(logger, tracer)
	// Add your GRPC options here

	grpcServer := grpc.NewGRPCServer(endpoints, options)

	var grpcListener net.Listener
	var err error

	g.Add(func() error {
		logger.Log("transport", "gRPC", "addr", *grpcAddr)
		grpcListener, err = net.Listen("tcp", *grpcAddr)
		if err != nil {
			return err
		}
		baseServer := grpc1.NewServer()
		pb.RegisterHelloServer(baseServer, grpcServer)
		gokit_foundation.RegistergRPCHealthSrv(baseServer)
		return baseServer.Serve(grpcListener)
	}, func(error) {
		if grpcListener != nil {
			grpcListener.Close()
		}
	})
}
func getServiceMiddleware(logger log.Logger) (mw []service.Middleware) {
	mw = []service.Middleware{}
	mw = addDefaultServiceMiddleware(logger, mw)
	// Append your middleware here

	return
}
func getEndpointMiddleware(logger log.Logger) (mw map[string][]endpoint1.Middleware) {
	mw = map[string][]endpoint1.Middleware{}
	duration := prometheus.NewSummaryFrom(prometheus1.SummaryOpts{
		Help:      "Request duration in seconds.",
		Name:      "request_duration_seconds",
		Namespace: "example",
		Subsystem: "hello",
	}, []string{"method", "success"})
	addDefaultEndpointMiddleware(logger, duration, mw)
	// Add you endpoint middleware here
	return
}
func initMetricsEndpoint(g *group.Group) {
	http1.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	var debugListener net.Listener
	var err error

	g.Add(func() error {
		logger.Log("transport", "debug/HTTP", "addr", *debugAddr)
		debugListener, err = net.Listen("tcp", *debugAddr)
		if err != nil {
			return err
		}
		return http1.Serve(debugListener, http1.DefaultServeMux)
	}, func(error) {
		if debugListener != nil {
			debugListener.Close()
		}
	})
}
func initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}
