package service

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/godcong/go-ffmpeg/oss"
	"github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// HandleFunc ...
type HandleFunc func(name, key string) error

var globalCancel context.CancelFunc

// Push ...
func Push(v *Streamer) {
	queue.RPush("queue", v.JSON())
}

// Pop ...
func Pop() *Streamer {
	pop := queue.LPop("queue").Val()
	return ParseStreamer(pop)
}

// StartQueue ...
func StartQueue(ctx context.Context, process int) {
	if queue == nil {
		queue = redis.NewClient(&redis.Options{
			Addr:     "",
			Password: "",              // no password set
			DB:       RedisQueueIndex, // use default DB
		})
	}
	var c context.Context
	c, globalCancel = context.WithCancel(ctx)
	//run with a new go channel
	go func() {
		threads := make(chan string, process)

		for i := 0; i < process; i++ {
			log.Println("start", i)
			go transferNothing(threads)
		}

		isStop := false
		count := 0
		for {
			select {
			case v := <-threads:
				log.Println("success:", v)
				if isStop {
					count++
					if count == process {
						return
					}
					continue
				}
				if s := Pop(); s != nil {
					go transfer(threads, s)
				} else {
					time.Sleep(3 * time.Second)
					go transferNothing(threads)
				}
				time.Sleep(5 * time.Second)
			case <-c.Done():
				isStop = true
				//default:
				//	log.Println("default")
				//	time.Sleep(3 * time.Second)
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

func transfer(chanints chan<- string, info *Streamer) {
	var err error
	server := oss.Server2()

	//fileSource := "./upload/" + util.GenerateRandomString(64)
	p := oss.NewProgress()
	p.SetObjectKey(info.ObjectKey)

	//fileName := filepath.Split(key)
	p.SetPath("./upload/")
	err = server.Download(p, info.FileName)

	if err != nil {
		log.Println()
		return
	}

	if info.Encrypt() {
		_ = info.KeyFile()
		err = ToM3U8WithKey(info.FileName)
	} else {
		err = ToM3U8(info.FileName)
	}

	if err != nil {
		//err = rdsQueue.Set(info.FileName, StatusFileWrong, 0).Err()
		if err != nil {
			log.Println(err)
		}
		return
	}
	log.Println("transferred:", *info)

	//err = rdsQueue.Set(info.FileName, StatusFinished, 0).Err()
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := http.PostForm("http://127.0.0.1:7790/v1/commit", url.Values{
		"id": []string{info.FileName},
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
			"id":       []string{info.FileName},
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

	chanints <- info.FileName
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
