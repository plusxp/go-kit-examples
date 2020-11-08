package crontask

import (
	"github.com/robfig/cron/v3"
	"go-util/_util"
	"hello/db"
	"math/rand"
	"time"
)

var cronMain = cron.New(cron.WithSeconds())

func Init() {
	// 定时任务开始前sleep一个随机值，避免多个进程读到相同数据，提高效率
	rand.Seed(time.Now().UnixNano())

	// 定时任务：打印程序资源消耗统计
	// 每2s执行
	_, err := cronMain.AddFunc("*/2 * * * * ?", printStatis)
	_util.PanicIfErr(err, nil)

	// 定时任务：假设场景，结算今日活动，并且给满足条件的人发送奖励（使用分布式锁避免多个进程同时执行该任务，造成多发奖励）
	// 分布式锁可以适用大部分此类场景，也有严谨度更高的方案：将任务状态持久化(不限存储源)
	// 每min的0s执行
	_, err = cronMain.AddFunc("0 * * * * ?", endupTodayActivity(db.GRedisDao.NewLock()))
	_util.PanicIfErr(err, nil)

	cronMain.Start()
}

func Stop() {
	cronMain.Stop()
}
