module gokit_foundation

go 1.12

require (
	github.com/go-kit/kit v0.10.0
	github.com/gorilla/mux v1.7.3
	github.com/hashicorp/consul/api v1.7.0
	github.com/sirupsen/logrus v1.4.2
	go-util v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.32.0
)

replace go-util => ../go-util
