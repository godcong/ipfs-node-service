package service

import (
	"context"
	"github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
		return
	}

	resp, err := http.PostForm("http://127.0.0.1:7790/v1/commit", url.Values{
		"id": []string{info.fileName},
		//"ipns": []string{uuid.NewV1().String()},
	})
	bytes, err := ioutil.ReadAll(resp.Body)
	log.Println(string(bytes), err)
	if err == nil {
		var cr CommitResult
		err := jsoniter.Unmarshal(bytes, &cr)
		if err != nil {
			log.Println(err)
			return
		}
		response, err := http.PostForm("http://127.0.0.1:7788/v0/ipfs/callback", url.Values{
			"id":       []string{info.fileName},
			"ipfs":     []string{cr.Detail.IpfsInfo.Hash},
			"ipns":     []string{cr.Detail.Ipns},
			"ipns_key": []string{cr.Detail.IpnsKey},
		})
		if err != nil {
			log.Println(err)
			return
		}
		by, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(by))
	}

	chanints <- info.fileName
}

// CommitResult ...
type CommitResult struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Detail struct {
		FileID   string `json:"fileID"`
		IpfsInfo struct {
			Hash string `json:"Hash"`
			Name string `json:"Name"`
			Size string `json:"Size"`
		} `json:"ipfsInfo"`
		Ipns     string `json:"ipns"`
		IpnsInfo struct {
			Name  string `json:"Name"`
			Value string `json:"Value"`
		} `json:"ipnsInfo"`
		IpnsKey string `json:"ipnsKey"`
	} `json:"detail"`
}

func transferNothing(threads chan<- string) {
	threads <- "nothing"
}
