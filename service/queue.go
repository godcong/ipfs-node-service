package service

import (
	"github.com/go-redis/redis"
	"log"
	"time"
)

// HandleFunc ...
type HandleFunc func(name, key string) error

var queue = NewStreamQueue()
var flag = false
var client = initClient()

func initClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})

	pong, err := client.Ping().Result()
	log.Println(pong)
	if err != nil {
		panic(err)
	}
	return client
}

// Push ...
func Push(v *StreamInfo) {
	queue.Push(v)
}

// Pop ...
func Pop() *StreamInfo {
	if !queue.IsEmpty() {
		return queue.Pop()
	}
	log.Println("nothing pop")
	return nil
}

// StartQueue ...
func StartQueue(process int) {
	flag = false
	//run with a new go channel
	go func() {
		threads := make(chan string, process)

		for i := 0; i < process; i++ {
			if flag {
				close(threads)
				return
			}
			log.Println("start", i)
			if s := Pop(); s != nil {
				go transfer(threads, s)
			} else {
				go transferNothing(threads)
			}
		}

		for {
			select {
			case v := <-threads:
				log.Println("success:", v)
				for {
					if flag {
						close(threads)
						return
					}
					if s := Pop(); s != nil {
						go transfer(threads, s)
						break
					}
					time.Sleep(5 * time.Second)
				}
			default:
				log.Println("default")
				time.Sleep(3 * time.Second)
				if flag {
					close(threads)
					return
				}
			}
		}
	}()
}

// StopQueue ...
func StopQueue() {
	//_ = q.db.Close()
	flag = true
}

func transfer(chanints chan<- string, info *StreamInfo) {
	key := info.KeyFile()
	_ = ToM3U8WithKey("./upload/", "./transfer/", info.fileName, key)
	time.Sleep(15 * time.Second)
	log.Println("transfer:", *info)
	//d, _ := json.Marshal(info)
	err := client.Set(info.fileName, info.key, 0).Err()
	if err != nil {
		log.Println(err)
	}

	chanints <- info.fileName
}

func transferNothing(chanints chan<- string) {
	log.Println("transferNothing")
	chanints <- "nothing"
}
