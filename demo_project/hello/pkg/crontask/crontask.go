package crontask

import (
	"github.com/leigg-go/go-util/_lock"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"log"
	"runtime"
	"time"
)

func printStatis() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	hostInfo, _ := host.Info()
	utilities, _ := cpu.Percent(time.Duration(time.Second*2), true)
	var cpuUtility float64 // 这个东西貌似不太准确
	for _, ut := range utilities {
		cpuUtility += ut
	}
	log.Printf("[crontask] hostInfo -- hostname:%s uptime:%d procs:%d hostid:%s, cpu-utility:%.2f", hostInfo.Hostname, hostInfo.Uptime, hostInfo.Procs, hostInfo.HostID, cpuUtility)
}

func endupTodayActivity(lock _lock.DistributedLock) func() {
	return func() {
		acquire := func() bool {
			// 此任务的内容决定它在同一时间只能运行一次，否则会发送多次奖励，这里使用redis分布式锁实现
			ok, err := lock.Lock()
			if err != nil {
				log.Printf("lock err:%v", err)
			}
			return ok
		}

		release := func() {
			if err := lock.UnLock(); err != nil {
				log.Printf("unlock err:%v", err)
			}
		}

		if !acquire() {
			log.Println("lock missed!")
			return
		}

		defer func() {
			// 一次性任务应该适当延长耗时，避免分布式环境中不同进程下某些原因导致的时钟误差问题从而导致第二个进程也能获得锁并执行任务
			time.Sleep(time.Second * 2)
			release()
		}()
		// 结算... 忽略过程

		// 发放奖励
		log.Print("发放奖励(注意应开启多进程测试，同一时刻应只有一个进程执行)。。")
	}
}
