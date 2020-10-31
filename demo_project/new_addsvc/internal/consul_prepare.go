package internal

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	stdconsul "github.com/hashicorp/consul/api"
	"gokit_foundation"
)

/*
服务注册到consul
*/

// protocol-svc_name-addr, e.g. grpc-UserServer-127.0.0.1:8888
const consulSvcIDFormat = "%s-%s-%s:%d"

func SvcRegisterTask(_ context.Context, logger log.Logger, svcName, svcHost string, port int) func(_ context.Context) error {
	return func(_ context.Context) error {
		logger.Log("NewSafeAsyncTask", "SvcRegister")
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
				GRPC:                           fmt.Sprintf("%s:%d/%s", svcHost, port, "grpc_health"), // or `grpc.health.v1.Health`
				Timeout:                        "5s",
				Interval:                       "5s",
				DeregisterCriticalServiceAfter: "15s", //check失败后多久删除本服务（位于consul中的服务条目）
			},
		}

		err := gokit_foundation.RegisterWithConsul(reg)
		return err
	}
}
