package main

import (
	"flag"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"gokit_foundation/gateway"
	"net/http"
	"os"
)

type MyGateWay struct {
	gateway.ExposedGateway
}

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8000", "Address for HTTP (JSON) server")
	)

	r := mux.NewRouter()

	root := gateway.New(r, *httpAddr, log.NewLogfmtLogger(os.Stderr))
	gw := &MyGateWay{root}

	{
		// declare a helloSvcRoute, then register some handler under the helloSvcRoute.
		helloSvcRoute := r.PathPrefix("/hello")

		helloSvcRoute.Methods("GET").Path("/sayhi/{name}").Handler(http.HandlerFunc(gw.SayHi))
	}

	// You may need process the err, or process err in foundation pkg in a unified way.
	_ = gw.Run()
}
