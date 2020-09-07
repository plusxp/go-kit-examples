package internal

import (
	"context"
	"fmt"
	stdconsul "github.com/hashicorp/consul/api"
	"go-util/_go"
	"gokit_foundation"
)

/*
服务注册到consul
*/

// protocol-svc_name-addr, e.g. grpc-UserServer-127.0.0.1:8888
const consulSvcIDFormat = "%s-%s-%s:%d"

func regTask(_ context.Context, svcName, svcHost string, port int) _go.AsyncTask {
	return func(context.Context, _go.Setter) {
		// consul agent配置，根据实际的填写
		reg := &stdconsul.AgentServiceRegistration{
			ID:                fmt.Sprintf(consulSvcIDFormat, "grpc", svcName, svcHost, port),
			Name:              svcName,
			Tags:              []string{"test", "gokit_svc"},
			Port:              port,
			Address:           svcHost,
			EnableTagOverride: false,
			// 配置实例本身的健康检查
			Check: &stdconsul.AgentServiceCheck{
				// consul服务必须能够访问这个地址，否则认为实例掉线
				GRPC:                           fmt.Sprintf("%s:%d/%s", svcHost, port, "grpc_health"),
				Timeout:                        "5s",
				Interval:                       "5s",
				DeregisterCriticalServiceAfter: "15s", //check失败后多久删除本服务
			},
		}

		gokit_foundation.RegisterWithConsul(reg)
	}
}

func SvcRegisterTask(ctx context.Context, svcName, svcHost string, port int) _go.AsyncTask {
	return regTask(ctx, svcName, svcHost, port)
}
