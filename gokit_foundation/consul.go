package gokit_foundation

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	stdconsul "github.com/hashicorp/consul/api"
	"net/http"
	"os"
	"time"
)

var DefaultRegister *consul.Registrar

func RegisterWithConsul(svcRegistration *stdconsul.AgentServiceRegistration) *consul.Registrar {
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
		panic(fmt.Sprintf("%s stdconsul.NewClient %v", time.Now().String()[:19], err))
	}

	kitConsulClient := consul.NewClient(consulClient)
	logger := log.NewLogfmtLogger(os.Stderr)

	registrar := consul.NewRegistrar(kitConsulClient, svcRegistration, log.With(logger, "component", "register"))
	registrar.Register()
	if DefaultRegister == nil {
		DefaultRegister = registrar
	}
	return registrar
}

func ConsulDeregister() {
	if DefaultRegister != nil {
		DefaultRegister.Deregister()
	}
}
