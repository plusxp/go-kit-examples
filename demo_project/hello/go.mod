module hello

go 1.12

require (
	github.com/go-kit/kit v0.10.0
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/leigg-go/go-util v0.0.3
	github.com/lightstep/lightstep-tracer-go v0.21.0
	github.com/oklog/oklog v0.3.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.4
	github.com/prometheus/client_golang v1.7.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/shirou/gopsutil v2.20.1+incompatible
	github.com/sony/gobreaker v0.4.1
	go-util v0.0.0-00010101000000-000000000000
	gokit_foundation v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20200904194848-62affa334b73
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/grpc v1.32.0
	sourcegraph.com/sourcegraph/appdash v0.0.0-20190731080439-ebfcffb1b5c0
)

replace (
	go-util => ../../go-util
	gokit_foundation => ../../gokit_foundation
)
