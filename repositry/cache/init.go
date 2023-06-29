package cache

import (
	"easy-drive/conf"
	"fmt"
	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func InitRedis() {
	rConfig := conf.Conf.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rConfig.RedisHost, rConfig.RedisPort),
		Password: rConfig.RedisPassword,
		DB:       rConfig.RedisDbName,
	})
	if _, err := client.Ping().Result(); err != nil {
		panic(err)
	}
	RedisClient = client
}
