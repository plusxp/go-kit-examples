package main

import (
	"flag"
	"github.com/gorilla/mux"
	"gokit_foundation"
	"gokit_foundation/gateway"
)

type MyGateWay struct {
	*gateway.Gateway
}

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8000", "Address for HTTP (JSON) server")
	)
	/*
		这里使用 https://github.com/gorilla/mux 作为路由器
	*/
	r := mux.NewRouter()

	root := gateway.New(r, *httpAddr, gokit_foundation.NewLogger())
	gw := &MyGateWay{root}

	{
		// 声明一个包含path前缀的子路由器
		helloSvcRoute := r.PathPrefix("/hello").Subrouter()

		// 为hellosvc这个服务router使用跨域策略
		// 注意：这比在全局使用跨域策略来的更灵活，如果没有使用Subrouter()，那就是在全局生效
		helloSvcRoute.Use(mux.CORSMethodMiddleware(helloSvcRoute))
		helloSvcRoute.Use(gateway.CORSHelper)

		// 用子路由器来注册更多的路由，他们共享主路由器的配置，注意：主路由器已有的配置不能在子路由器中再次配置，可能会导致一些问题
		helloSvcRoute.HandleFunc("/sayhi/{name}", gw.SayHi).Methods("GET", "OPTIONS")
		helloSvcRoute.HandleFunc("/make_a_date/{date:\\d{4}-\\d\\d-\\d\\d}", gw.SayHi).Methods("GET", "OPTIONS")

	}

	// 直接运行！
	_ = gw.Run()
}
