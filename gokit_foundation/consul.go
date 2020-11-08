package gokit_foundation

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	stdconsul "github.com/hashicorp/consul/api"
	"go-util/_util"
	"net/http"
	"os"
	"time"
)

var DefaultRegister *consul.Registrar

// protocol-svc_name-addr, e.g. grpc-UserServer-127.0.0.1:8888
const consulSvcIDFormat = "%s-%s-%s:%d"

func MustRegisterSvc(svcName, svcHost string, port int, tags []string) {
	// consul agent配置，根据实际的填写
	tags = append(tags, "gokit_svc")
	reg := &stdconsul.AgentServiceRegistration{
		ID:                fmt.Sprintf(consulSvcIDFormat, "grpc", svcName, svcHost, port),
		Name:              svcName,
		Tags:              tags,
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
	err := RegisterWithConsul(reg)
	_util.PanicIfErr(err, nil)
}

func RegisterWithConsul(svcRegistration *stdconsul.AgentServiceRegistration) error {
	if DefaultRegister != nil {
		return nil
	}
	consulAddr := os.Getenv("CONSUL_ADDR")
	if consulAddr == "" {
		//panic(fmt.Sprintf("%s CONSUL_ADDR not set", time.Now().String()[:19]))
		consulAddr = "127.0.0.1:8500"
	}
	// 这个client是针对consul，不是服务
	consulClient, err := stdconsul.NewClient(&stdconsul.Config{
		Address: consulAddr,
		HttpClient: &http.Client{
			Timeout: time.Second * 2,
		},
		Scheme: "http", // default
	})
	if err != nil {
		return err
	}

	kitConsulClient := consul.NewClient(consulClient)
	logger := log.NewLogfmtLogger(os.Stderr)

	registrar := consul.NewRegistrar(kitConsulClient, svcRegistration, log.With(logger, "component", "register"))
	registrar.Register()
	DefaultRegister = registrar
	return nil
}

func ConsulDeregister() {
	if DefaultRegister != nil {
		DefaultRegister.Deregister()
	}
}
