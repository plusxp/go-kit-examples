package service

import (
	"github.com/leigg-go/go-util/_redis"
	"hello/config"
	"hello/pkg/crontask"
)

// 短时间的初始化任务，这种不能用g.Add
func initShortTimeTask() {
	_redis.MustInitDefClient(config.GetRedisConf())
	crontask.Init()
}
func onClose() {
	crontask.Stop()
}
