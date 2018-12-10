package service

import (
	"sync"
	"time"
)

var queue sync.Pool

func Processor() {
	for {
		if v := queue.Get(); v != nil {

		}
		time.Sleep(5 * time.Second)
	}
}
