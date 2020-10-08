module gokit_foundation

go 1.12

require (
	github.com/go-kit/kit v0.10.0
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.4.1
	github.com/gorilla/mux v1.7.3
	github.com/hashicorp/consul/api v1.7.0
	go-util v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/jose.v1 v1.0.0-20161127122323-a941c3995164
)

replace go-util => ../go-util
