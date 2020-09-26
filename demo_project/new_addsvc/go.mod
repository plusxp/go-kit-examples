module new_addsvc

go 1.12

require (
	github.com/go-kit/kit v0.10.0
	github.com/golang/protobuf v1.4.1
	github.com/google/go-cmp v0.5.1 // indirect
	github.com/hashicorp/consul/api v1.7.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.3.0
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/sony/gobreaker v0.4.1
	github.com/stretchr/testify v1.6.1 // indirect
	go-util v0.0.0-00010101000000-000000000000
	gokit_foundation v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sys v0.0.0-20200810151505-1b9f1253b3ed // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace go-util => ../../go-util

replace gokit_foundation => ../../gokit_foundation
