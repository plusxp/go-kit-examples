package db

import (
	"github.com/leigg-go/go-util/_lock"
	"time"
)

func (rds *RedisType) NewLock() _lock.DistributedLock {
	opt := _lock.LockOption{Timeout: time.Second, Retry: false}
	expire := time.Second * 2 // key过期时间，略大于任务耗时+网络传输耗时即可
	return _lock.NewDistributedLockByRedis(rds.cli, "lock_key_endupTodayActivity", nil, expire, opt)
}

// 所有redis的操作都应该定义以RedisType为receiver的Method
//
//func (rds *RedisType) GetUserInfoFromCache() (string, error) {
//	return rds.cli.Get("uid_xxx").Result()
//}
