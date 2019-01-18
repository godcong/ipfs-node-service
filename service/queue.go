package service

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
	"log"
	"strconv"
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
			go transferNothing(threads, strconv.Itoa(i))
		}
		isStop := false
		count := 0
		for {
			select {
			case v := <-threads:
				if isStop {
					count++
					if count == process {
						log.Println("stoped")
						return
					}
					continue
				}
				if s := Pop(); s != nil {
					go transfer(threads, s, v)
				} else {
					go transferNothing(threads, v)
				}
			case <-c.Done():
				isStop = true
			default:
				time.Sleep(1 * time.Second)
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

func transfer(ch chan<- string, info *Streamer, idx string) {
	var err error
	chanRes := idx
	defer func() {
		if err != nil {
			log.Println(err)
			chanRes = idx + info.FileName() + err.Error()
		}
		ch <- chanRes
	}()

	queue.Set(info.ID, StatusDownloading, 0)
	err = download(info)
	if err != nil {
		return
	}

	queue.Set(info.ID, StatusTransferring, 0)
	if info.Encrypt() {
		_ = info.KeyFile()
		err = toM3U8WithKey(info.ID, info.SourceFile(), info.FileDest, info.KeyInfoName)
	} else {
		err = toM3U8(info.ID, info.SourceFile(), info.FileDest)
	}

	if err != nil {
		return
	}

	detail, err := commitToIPNS(info.ID, info.DestPath())
	if err != nil {
		return
	}

	var qr QueueResult

	err = mapstructure.Decode(detail, &qr)
	if err != nil {
		return
	}

	//stream.Callback = Config().Callback.Type

	err = NewBack().Callback(&qr)

	//response, err := http.PostForm("http://127.0.0.1:7788/v0/ipfs/callback", url.Values{
	//	"id":       []string{info.FileName()},
	//	"ipfs":     []string{cr.Detail.IpfsInfo.Hash},
	//	"ipns":     []string{cr.Detail.Ipns},
	//	"ipns_key": []string{cr.Detail.IpnsKey},
	//})
	if err != nil {
		return
	}
	//by, err := ioutil.ReadAll(response.Body)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//log.Println(string(by))

}

// QueueResult ...
type QueueResult struct {
	ID     string `mapstructure:"id"`
	FSInfo struct {
		Hash string `mapstructure:"hash"`
		Name string `mapstructure:"name"`
		Size string `mapstructure:"size"`
	} `mapstructure:"fs_info"`
	NSInfo struct {
		Name  string `mapstructure:"name"`
		Value string `mapstructure:"value"`
	} `mapstructure:"ns_info"`
}

// JSON ...
func (r *QueueResult) JSON() string {
	s, _ := jsoniter.MarshalToString(r)
	return s
}

func transferNothing(threads chan<- string, idx string) {
	time.Sleep(3 * time.Second)
	threads <- idx
}
