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
	*badger.DB
	queue sync.Pool
}

var queue Queue

//InitDB ...
func InitDB() {
	var err error
	options := badger.DefaultOptions
	options.Dir = "/tmp/badger"
	options.ValueDir = "/tmp/badger"
	queue.DB, err = badger.Open(options)
	if err != nil {
		panic(err)
	}
}

// DB ...
func DB() *badger.DB {
	return queue.DB
}

// Add ...
func Add(v *Streamer) {
	queue.queue.Put(v)
	return
}
