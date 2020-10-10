package gokit_foundation

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	stdconsul "github.com/hashicorp/consul/api"
	"net/http"
	"os"
	"time"
)

var DefaultRegister *consul.Registrar

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
			Timeout: time.Second * 3,
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
