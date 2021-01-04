module gokit_foundation

go 1.12

require (
	github.com/HdrHistogram/hdrhistogram-go v1.0.1 // indirect
	github.com/go-kit/kit v0.10.0
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.4.1
	github.com/gorilla/mux v1.7.3
	github.com/hashicorp/consul/api v1.7.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	go-util v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/jose.v1 v1.0.0-20161127122323-a941c3995164
)

replace go-util => ../go-util
