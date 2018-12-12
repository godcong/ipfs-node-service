package service

import (
	"github.com/dgraph-io/badger"
	"sync"
)

// Streamer ...
type Streamer struct {
	key      string
	fileName string
	status   int
	step     int
}

// NewStreamer ...
func NewStreamer(key string, fileName string) *Streamer {
	return &Streamer{key: key, fileName: fileName}
}

// Status ...
func (s *Streamer) Status() int {
	return s.status
}

// SetStatus ...
func (s *Streamer) SetStatus(status int) {
	s.status = status
}

// Step ...
func (s *Streamer) Step() int {
	return s.step
}

// SetStep ...
func (s *Streamer) SetStep(step int) {
	s.step = step
}

// FileName ...
func (s *Streamer) FileName() string {
	return s.fileName
}

// SetFileName ...
func (s *Streamer) SetFileName(fileName string) {
	s.fileName = fileName
}

// Key ...
func (s *Streamer) Key() string {
	return s.key
}

// SetKey ...
func (s *Streamer) SetKey(key string) {
	s.key = key
}

// Queue ...
type Queue struct {
	db    *badger.DB
	queue sync.Pool
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
func (q *Queue) Push(v *Streamer) {
	q.queue.Put(v)
}

// Pop ...
func (q *Queue) Pop() *Streamer {
	if v := q.queue.Get(); v != nil {
		return v.(*Streamer)
	}
	return nil
}
