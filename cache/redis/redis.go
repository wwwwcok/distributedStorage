package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

var RedisCli *redis.Client

func init() {
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.0.90:6379",
		Password: "",
		DB:       0,
	})

	RedisCli = client

	err := client.Set("test", "测试redis连接", 0).Err()
	if err != nil {
		fmt.Println("init测试redis连接失败:", err)
	}

}
