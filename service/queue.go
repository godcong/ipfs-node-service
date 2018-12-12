package service

import (
	"github.com/go-redis/redis"
	"github.com/godcong/go-ffmpeg/openssl"
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
	uri      string
	src      string
	dst      string
}

func (s *StreamInfo) Dst() string {
	return s.dst
}

func (s *StreamInfo) SetDst(dst string) {
	s.dst = dst
}

func (s *StreamInfo) Src() string {
	return s.src
}

func (s *StreamInfo) SetSrc(src string) {
	s.src = src
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

func (s *StreamInfo) Uri() string {
	return s.uri
}

func (s *StreamInfo) SetUri(uri string) {
	s.uri = uri
}

// KeyFile ...
func (s *StreamInfo) KeyFile() string {
	err := openssl.KeyFile("./transfer/", s.fileName, s.key, "localhost:8080/stream", true)
	if err != nil {
		return ""
	}
	return "./transfer/" + s.fileName + "_keyfile"
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
		DB:       1,  // use default DB
	})

	pong, err := client.Ping().Result()
	log.Println(pong)
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
	//if q.client == nil {
	//	q.client = initClient()
	//}
	return q.client
}

// Push ...
func (q *Queue) Push(v *StreamInfo) {
	log.Println("pushing", v)
	q.queue.Put(v)
}

// Pop ...
func (q *Queue) Pop() *StreamInfo {

	if v := q.queue.Get(); v != nil {
		log.Println("poping", v)
		return v.(*StreamInfo)
	}
	log.Println("poping", "nil")
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
				log.Println("default")
				time.Sleep(5 * time.Second)
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
	//key := info.KeyFile()
	//_ = ToM3U8WithKey("./upload/", "./transfer/", info.fileName, key)
	time.Sleep(10 * time.Second)
	log.Println("transfer:", *info)
	//d, _ := json.Marshal(info)
	err := queue.Client().Set(info.fileName, info.key, 0).Err()
	if err != nil {
		log.Println(err)
	}

	chanints <- 1
}
