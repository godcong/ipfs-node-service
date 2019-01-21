package service

import (
	"fmt"
	"github.com/godcong/node-service/config"
	"github.com/godcong/node-service/oss"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// RunMain 主线程
func RunMain() {
	log.Println("run main")
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cfg := config.Config()

	oss.InitOSS(cfg)

	//rest start
	rest := NewRestServer(cfg)
	_ = Router(rest.Engine)
	rest.Start()

	//grpc start
	grpc := NewGRPCServer(cfg)
	grpc.Start()

	queue := NewQueueServer(cfg)
	queue.Processes = 5
	queue.Start()

	go func() {
		sig := <-sigs
		//bm.Stop()
		fmt.Println(sig, "exiting")
		rest.Stop()
		grpc.Stop()
		queue.Stop()
		done <- true
	}()
	<-done
}

// NewBack ...
func NewBack() StreamerCallback {
	cfg := config.Config()
	if cfg != nil && cfg.Callback.Type == "grpc" {
		return NewGRPCBack(cfg)
	}
	return NewRestBack(cfg)
}
