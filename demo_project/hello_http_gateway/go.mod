module hello_http_gateway

go 1.12

require (
	github.com/gorilla/mux v1.8.0
	gokit_foundation v0.0.0-00010101000000-000000000000
	hello v0.0.0-00010101000000-000000000000
)

replace (
	go-util => ../../go-util
	gokit_foundation => ../../gokit_foundation
	hello => ../hello
)
