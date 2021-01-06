module hello

go 1.15

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/go-kit/kit v0.10.0
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/hashicorp/consul/api v1.7.0
	github.com/leigg-go/go-util v0.0.4
	github.com/oklog/oklog v0.3.2
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.7.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/shirou/gopsutil v2.20.1+incompatible
	github.com/sony/gobreaker v0.4.1
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	go-util v0.0.0-00010101000000-000000000000
	gokit_foundation v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20200904194848-62affa334b73
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/grpc v1.32.0
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.20.6
)

replace (
	go-util => ../../go-util
	gokit_foundation => ../../gokit_foundation
)
