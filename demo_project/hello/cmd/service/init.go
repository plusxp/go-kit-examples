package service

import (
	"gokit_foundation"
	"hello/config"
	"hello/db"
	"hello/pkg/crontask"
	"time"
)

func initFirstly(logger *gokit_foundation.Logger, grpcHost string, grpcPort int) {
	db.Init(logger) // 先启动db
	crontask.Init()
	// 等一下，待所有内部服务准备就绪后再上线服务
	// 1s表示所有内部服务需在1s内进入就绪状态，任何影响服务正常运行的错误都应该立即报错， 时间可调整
	time.AfterFunc(time.Second, func() {
		gokit_foundation.MustRegisterSvc(config.SvcName, grpcHost, grpcPort, []string{"test"})
	})
}

func onClose() {
	gokit_foundation.ConsulDeregister() // 先下线
	crontask.Stop()
	db.Close() // 最后停止db
}
