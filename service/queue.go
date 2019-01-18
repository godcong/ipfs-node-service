package service

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
	"log"
	"time"
)

// HandleFunc ...
type HandleFunc func(name, key string) error
type QueueServer struct {
	*redis.Client
	Processes int
	cancel    context.CancelFunc
}

var queue *QueueServer

// Push ...
func Push(v *Streamer) {
	queue.RPush("node_queue", v.JSON())
}

// Pop ...
func Pop() *Streamer {
	pop := queue.LPop("node_queue").Val()
	return ParseStreamer(pop)
}

func transfer(ch chan<- string, info *Streamer) {
	var err error
	chanRes := info.FileName()
	defer func() {
		if err != nil {
			log.Println(err)
			chanRes = info.FileName() + err.Error()
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

func transferNothing(threads chan<- string) {
	time.Sleep(3 * time.Second)
	threads <- ""
}

func NewQueueServer() *QueueServer {
	client := redis.NewClient(&redis.Options{
		Addr:     DefaultString(config.Queue.HostPort, ":6379"),
		Password: DefaultString(config.Queue.Password, ""), // no password set
		DB:       config.Queue.DB,                          // use default DB
	})
	return &QueueServer{
		Client: client,
	}
}

func (s *QueueServer) Start() {
	pong, err := s.Ping().Result()
	if err != nil {
		panic(err)
	}
	log.Println(pong)

	var c context.Context
	c, s.cancel = context.WithCancel(context.Background())
	//run with a new go channel
	go func() {
		threads := make(chan string, s.Processes)

		for i := 0; i < s.Processes; i++ {
			log.Println("start", i)
			go transferNothing(threads)
		}
		//isStop := false
		//count := 0
		for {
			select {
			case v := <-threads:
				if v != "" {
					log.Println("success: ", v)
				}
				//if isStop {
				//	count++
				//	if count ==  s.Processes {
				//		log.Println("stopped")
				//		return
				//	}
				//	continue
				//}
				if s := Pop(); s != nil {
					go transfer(threads, s)
				} else {
					go transferNothing(threads)
				}
			case <-c.Done():
				//isStop = true
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func (s *QueueServer) Stop() {
	if s.cancel == nil {
		return
	}
	s.cancel()
}
