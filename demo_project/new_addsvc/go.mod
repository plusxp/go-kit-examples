module new_addsvc

go 1.12

require (
	github.com/HdrHistogram/hdrhistogram-go v1.0.1 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/go-kit/kit v0.10.0
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/golang/protobuf v1.4.1
	github.com/hashicorp/consul/api v1.7.0
	github.com/leigg-go/go-util v0.0.4
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.3.0
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/shirou/gopsutil v2.20.9+incompatible
	github.com/sony/gobreaker v0.4.1
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	go-util v0.0.0-00010101000000-000000000000
	gokit_foundation v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sys v0.0.0-20200810151505-1b9f1253b3ed // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace go-util => ../../go-util

replace gokit_foundation => ../../gokit_foundation
