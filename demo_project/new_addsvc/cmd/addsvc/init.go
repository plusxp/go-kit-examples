package main

import (
	"github.com/leigg-go/go-util/_redis"
	"gokit_foundation"
	"new_addsvc/config"
	"new_addsvc/pkg/crontask"
	"time"
)

// 短时间的初始化任务，这种不能用g.Add
func initFirstly(grpcHost string, grpcPort int) {
	_redis.MustInitDef(config.GetRedisConf())
	crontask.Init()
	// 等一下，待所有内部服务准备就绪后再上线服务
	// 1s表示所有内部服务需在1s内进入就绪状态，任何影响服务正常运行的错误都应该立即报错， 时间可根据完成具体初始化任务所需的时间来调整
	time.AfterFunc(time.Second, func() {
		gokit_foundation.MustRegisterSvc(config.SvcName, grpcHost, grpcPort, []string{"test"})
	})
}
func onClose() {
	gokit_foundation.ConsulDeregister() // 先下线
	crontask.Stop()
	_redis.Close()
}
