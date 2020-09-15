package hello_http_gateway

import (
	"flag"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"gokit_foundation/gateway"
	"net/http"
)

type MyGateWay struct {
	gateway.ExposedGateway
}

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8000", "Address for HTTP (JSON) server")
	)

	r := mux.NewRouter()

	gw := gateway.New(r, *httpAddr, log.NewNopLogger())

	{
		// declare a helloSvcRoute, then register some handler under the helloSvcRoute.
		helloSvcRoute := r.PathPrefix("/hello")

		helloSvcRoute.Methods("GET").Path("/sayhi/{name}").Handler(http.HandlerFunc(SayHi))
	}

	// You may need process the err, or process err in foundation pkg in a unified way.
	_ = gw.Run()
}
