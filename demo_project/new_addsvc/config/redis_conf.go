package config

import (
	"github.com/go-redis/redis"
	"time"
)

func GetRedisConf() *redis.Options {
	var opt *redis.Options

	/*
		从配置文件/第三方支持kv存储的平台获取(etcd/consul/redis等)，使用redis还需要自行实现监听key的变化
		这里忽略读取过程...
	*/
	opt = &redis.Options{
		Addr:         "127.0.0.1:6379",
		Password:     "123",
		DB:           0,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		MinIdleConns: 1,
		IdleTimeout:  3 * time.Second,
	}

	return opt
}
