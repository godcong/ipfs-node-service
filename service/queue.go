package service

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/godcong/node-service/oss"
	"github.com/json-iterator/go"
	"github.com/satori/go.uuid"
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
	queue.RPush("node_queue", v.JSON())
}

// Pop ...
func Pop() *Streamer {
	pop := queue.LPop("node_queue").Val()
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
						log.Println("stoped")
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
			case <-c.Done():
				isStop = true
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

func download(info *Streamer) error {
	server := oss.Server2()
	queue.Set(info.ID, StatusDownloading, 0)
	p := oss.NewProgress()
	p.SetObjectKey(info.ObjectKey)
	p.SetPath(config.Media.Upload)
	err := server.Download(p, info.FileName())
	if err != nil {
		//chanRes = err.Error()
		log.Println(err)
		return err
	}
	return nil
}

func transfer(ch chan<- string, info *Streamer) {
	var err error
	chanRes := info.FileName()
	defer func() {
		ch <- chanRes
	}()

	if err != nil {
		chanRes = err.Error()
		log.Println(err)
		return
	}

	if info.Encrypt() {
		_ = info.KeyFile()
		err = ToM3U8WithKey(info.FileName())
	} else {
		err = ToM3U8(info.ID, info.FileName(), info.FileSource, info.FileDest)
	}

	if err != nil {
		if err != nil {
			chanRes = err.Error()
			log.Println(err)
		}
		return
	}
	log.Println("transferred:", *info)

	if err != nil {
		chanRes = err.Error()
		log.Println(err)
		return
	}
	ipfsInfo, err := api.AddDir(info.FileDest + "/" + info.FileName() + "/")
	if err != nil {
		chanRes = err.Error()
		log.Println(err)
		return
	}

	keyID := ""
	if ipns != "" {
		keyID, err = rdsIPNS.Get(ipns).Result()
	}
	log.Println(ipns, keyID, "error:", err)
	if ipns == "" || err != nil {
		keyID = uuid.NewV1().String()
		m, err := api.Key().Gen(keyID, "rsa", 2048)
		if err != nil {
			resultFail(ctx, err.Error())
			return
		}
		ipns = m["Id"]
		log.Println("ipns:", ipns, "key:", keyID)
		err = rdsIPNS.Set(ipns, keyID, 0).Err()
		if err != nil {
			log.Println(err)
			resultFail(ctx, err.Error())
			return
		}
	}
	log.Println(ipfsInfo)
	ipnsInfo, err := api.Name().PublishWithKey("/ipfs/"+ipfsInfo["Hash"], keyID)
	log.Println(ipnsInfo, err)
	if err != nil {
		chanRes = err.Error()

		return
	}
	log.Println(info.ID, ipnsInfo)
	//resultOK(ctx, gin.H{
	//	"fileID":   id,
	//	"ipns":     ipns,
	//	"ipnsKey":  keyID,
	//	"ipfsInfo": ipfsInfo,
	//	"ipnsInfo": ipnsInfo,
	//})
	resp, err := http.PostForm("http://127.0.0.1:7790/v1/commit", url.Values{
		"id": []string{info.FileName()},
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
			"id":       []string{info.FileName()},
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

	chanints <- info.FileName()
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
