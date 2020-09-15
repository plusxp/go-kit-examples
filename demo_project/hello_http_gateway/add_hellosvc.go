package hello_http_gateway

import (
	"context"
	"github.com/gorilla/mux"
	"hello/client/grpc"
	"net/http"
)

/*
API defined here.
*/

// restful-api
func (gw *MyGateWay) SayHi(w http.ResponseWriter, r *http.Request) {

	// Extract params from path
	v := mux.Vars(r)

	name := v["name"]
	if name == "" {

	}

	helloSvcClient := grpc.NewClient()
	rsp, err := helloSvcClient.SayHi(context.Background(), "Jack Ma")
	gw.JSON(w, rsp)
}
