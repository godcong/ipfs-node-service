package service

import (
	"context"
	"log"
	"time"
)

// HandleFunc ...
type HandleFunc func(name, key string) error

var queue = NewStreamQueue()
var globalCancel context.CancelFunc

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
func StartQueue(ctx context.Context, process int) {
	var c context.Context
	c, globalCancel = context.WithCancel(ctx)
	//run with a new go channel
	go func() {
		threads := make(chan string, process)

		for i := 0; i < process; i++ {
			log.Println("start", i)
			go transferNothing(threads)

		}

		for {
			select {
			case v := <-threads:
				log.Println("success:", v)
				if s := Pop(); s != nil {
					go transfer(threads, s)
				} else {
					time.Sleep(3 * time.Second)
					go transferNothing(threads)
				}
				time.Sleep(5 * time.Second)
			case <-c.Done():
				break
			default:
				log.Println("default")
				time.Sleep(3 * time.Second)
			}
		}
	}()
}

// StopQueue ...
func StopQueue() {
	if globalCancel == nil {
		return
	}
	globalCancel()
}

func transfer(chanints chan<- string, info *StreamInfo) {
	var err error
	if info.Encrypt() {
		_ = info.KeyFile()
		err = ToM3U8WithKey(info.fileName)
	} else {

		err = ToM3U8(info.fileName)
	}

	if err != nil {
		err = rdsQueue.Set(info.fileName, StatusFileWrong, 0).Err()
		if err != nil {
			log.Println(err)
		}
		return
	}
	log.Println("transferred:", *info)

	err = rdsQueue.Set(info.fileName, StatusFinished, 0).Err()
	if err != nil {
		log.Println(err)
	}

	chanints <- info.fileName
}

func transferNothing(chanints chan<- string) {
	log.Println("transferNothing")
	chanints <- "nothing"
}
