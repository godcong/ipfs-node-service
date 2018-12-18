package service

import (
	"github.com/go-redis/redis"
	"log"
)

var rdsQueue = newRedisWithDB(1)
var rdsIPNS = newRedisWithDB(2)

// newRedisWithDB ...
func newRedisWithDB(idx int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",  // no password set
		DB:       idx, // use default DB
	})

	pong, err := client.Ping().Result()
	log.Println("pong:", pong)
	if err != nil {
		panic(err)
	}
	return client
}
