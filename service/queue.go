package service

import (
	"github.com/go-redis/redis"
	"log"
	"sync"
	"time"
)

// HandleFunc ...
type HandleFunc func(name, key string) error

// StreamInfo ...
type StreamInfo struct {
	key      string
	fileName string
}

// NewStreamer ...
func NewStreamer(key string, fileName string) *StreamInfo {
	return &StreamInfo{key: key, fileName: fileName}
}

// FileName ...
func (s *StreamInfo) FileName() string {
	return s.fileName
}

// SetFileName ...
func (s *StreamInfo) SetFileName(fileName string) {
	s.fileName = fileName
}

// Key ...
func (s *StreamInfo) Key() string {
	return s.key
}

// SetKey ...
func (s *StreamInfo) SetKey(key string) {
	s.key = key
}

// Queue ...
type Queue struct {
	flag   bool
	client *redis.Client
	queue  sync.Pool
	handle HandleFunc
}

// NewQueue ...
func NewQueue(client *redis.Client) *Queue {
	return &Queue{client: client}
}

var queue *Queue

func initClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return client
}

// InitQueue ...
func InitQueue() *Queue {
	queue = NewQueue(initClient())
	return queue
}

// Client ...
func (q *Queue) Client() *redis.Client {
	if q.client == nil {
		q.client = initClient()
	}
	return q.client
}

// Push ...
func (q *Queue) Push(v *StreamInfo) {
	q.queue.Put(v)
}

// Pop ...
func (q *Queue) Pop() *StreamInfo {
	if v := q.queue.Get(); v != nil {
		return v.(*StreamInfo)
	}
	return nil
}

// Handle ...
func (q *Queue) Handle() HandleFunc {
	return q.handle
}

// SetHandle ...
func (q *Queue) SetHandle(handle HandleFunc) {
	q.handle = handle
}

// Start ...
func (q *Queue) Start(process int) {
	q.flag = false
	//run with a new go channel
	go func() {
		threads := make(chan int, process)

		for i := 0; i < process; i++ {
			for {
				if q.flag {
					close(threads)
					return
				}
				if s := q.Pop(); s != nil {
					log.Println("start", i)
					go transfer(threads, s)
					break
				}
				time.Sleep(5 * time.Second)
			}
		}

		for {
			select {
			case _ = <-threads:
				for {
					if q.flag {
						close(threads)
						return
					}
					if s := q.Pop(); s != nil {
						log.Println("thread run")
						go transfer(threads, s)
						break
					}
					time.Sleep(5 * time.Second)
				}
			default:
				if q.flag {
					close(threads)
					return
				}
			}
		}
	}()
}

// Stop ...
func (q *Queue) Stop() {
	//_ = q.db.Close()
	q.flag = true
}

func transfer(chanints chan<- int, info *StreamInfo) {
	_ = ToM3U8("./transfer", info.fileName)
	log.Println("transfer:", *info)
	chanints <- 1
}
