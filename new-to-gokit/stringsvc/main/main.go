package main

import (
	"github.com/chaseSpace/go-kit-examples/new-to-gokit/stringsvc"
	"net/http"
)

/*
一个最小的go kit服务
*/

func main() {

	http.Handle("/uppercase", stringsvc.UppercaseHandler)
	http.Handle("/count", stringsvc.CountHandler)
	stringsvc.Logger.Log(http.ListenAndServe(":8080", nil))
}
