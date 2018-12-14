package service

import (
	"github.com/go-redis/redis"
	"log"
)

var client = initClient()

func initClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})

	pong, err := client.Ping().Result()
	log.Println("pong:", pong)
	if err != nil {
		panic(err)
	}
	return client
}
