package service

import (
	"github.com/dgraph-io/badger"
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
	db     *badger.DB
	queue  sync.Pool
	handle HandleFunc
}

// NewQueue ...
func NewQueue(db *badger.DB) *Queue {
	return &Queue{db: db}
}

var queue *Queue

//InitDB ...
func initDB() *badger.DB {
	var err error
	options := badger.DefaultOptions
	options.Dir = "/home/badger"
	options.ValueDir = "/home/badger"
	db, err := badger.Open(options)
	if err != nil {
		panic(err)
	}
	return db
}

func init() {
	queue = NewQueue(initDB())
}

// InitQueue ...
func InitQueue() *Queue {
	return queue
}

// DB ...
func (q *Queue) DB() *badger.DB {
	if q.db == nil {
		q.db = initDB()
	}
	return q.db
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

// Run ...
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
	q.flag = true
}

func transfer(chanints chan<- int, info *StreamInfo) {

}
