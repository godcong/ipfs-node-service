package service

import (
	"github.com/go-redis/redis"
)

// RedisQueueIndex ...
const RedisQueueIndex = 1

//var rdsQueue = newRedisWithDB(1)
var rdsIPNS = newRedisWithDB(2)

// newRedisWithDB ...
func newRedisWithDB(idx int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",  // no password set
		DB:       idx, // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return client
}
