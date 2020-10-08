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

	lgr := gokit_foundation.NewLogger(nil)
	root := gateway.New(r, *httpAddr, lgr)

	gw := &MyGateWay{root}

	{
		// 声明一个包含path前缀的子路由器
		helloSvcRoute := r.PathPrefix("/hello").Subrouter()

		// 为helloRouter使用跨域策略
		// 注意：这比在全局使用跨域策略来的更灵活，如果没有使用Subrouter()，那就是在全局生效
		helloSvcRoute.Use(mux.CORSMethodMiddleware(helloSvcRoute))
		helloSvcRoute.Use(gateway.CORSHelper)

		// 用子路由器来注册更多的路由，他们共享主路由器的配置，注意：主路由器已有的配置不能在子路由器中再次配置，可能会导致一些问题
		helloSvcRoute.HandleFunc("/sayhi/{name}", gw.SayHi).Methods("GET", "OPTIONS")
		// 此正则匹配类似2020-10-10的日期参数，若匹配不上返回404
		helloSvcRoute.HandleFunc("/make_a_date/{date:\\d{4}-\\d\\d-\\d\\d}/{want_say:.*}", gw.MakeADate).Methods("GET", "OPTIONS")
		// 一个post接口
		helloSvcRoute.HandleFunc("/update_user_info", gw.UpdateUserInfo).Methods("POST", "OPTIONS")
	}

	// 直接运行！(先启动hello服务)
	_ = gw.Run()
}

/*
如何测试：
	curl http://127.0.0.1:8000/hello/sayhi/Hanmeimei
	# %20 在URL中表示空格
	curl http://127.0.0.1:8000/hello/make_a_date/2020-10-01/Do%20you%20willing%20to%20date%20with%20me?

	# 模拟一个服务端错误(12-12是暗号，会被服务端特殊处理)：
	curl http://127.0.0.1:8000/hello/make_a_date/2020-12-12

	update_user_info接口测试参看main_test.go

多次且快速的发送/hello/make_a_date/2020-12-12请求，断路器将会打开, 将会看到错误变更:
{"caller":"demo_project/gateway/hellosvc.go:77","err":"rpc error: code = Unknown desc = dependency svc err","ts":"2020-09-19 12:07:09"}
{"caller":"demo_project/gateway/hellosvc.go:77","err":"rpc error: code = Unknown desc = dependency svc err","ts":"2020-09-19 12:07:09"}
{"caller":"demo_project/gateway/hellosvc.go:77","err":"circuit breaker is open","ts":"2020-09-19 12:07:09"}
{"caller":"demo_project/gateway/hellosvc.go:77","err":"circuit breaker is open","ts":"2020-09-19 12:07:09"}
{"caller":"demo_project/gateway/hellosvc.go:77","err":"circuit breaker is open","ts":"2020-09-19 12:07:10"}
*/
