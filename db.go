package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	rdb *redis.Client
)

// 初始化连接
func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "172.16.1.3:6379",
		Password: "djs@12316", // no password set
		DB:       0,           // use default DB
		PoolSize: 10,          // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	return err
}
