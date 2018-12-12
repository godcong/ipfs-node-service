package service

import (
	"github.com/dgraph-io/badger"
	"sync"
)

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
func Add(v interface{}) {
	queue.queue.Put(v)
	return
}
