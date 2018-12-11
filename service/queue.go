package service

import (
	"github.com/dgraph-io/badger"
	"sync"
	"time"
)

var db *badger.DB
var queue sync.Pool

// Processor ...
func Processor() {
	for {
		if v := queue.Get(); v != nil {

		}
		time.Sleep(5 * time.Second)
	}
}

func InitDB() {
	var err error
	options := badger.DefaultOptions
	options.Dir = "/tmp/badger"
	options.ValueDir = "/tmp/badger"
	db, err = badger.Open(options)
	if err != nil {

	}
}
