package service

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var configPath = flag.String("path", "config.toml", "load config file from path")

// RunMain 主线程
func RunMain() {
	if !flag.Parsed() {
		flag.Parse()
	}
	err := Initialize(*configPath)
	if err != nil {
		panic(err)
	}
	log.Println("run main")
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//rest start
	rest := NewRestServer()
	_ = Router(rest.Engine)
	rest.Start()

	//grpc start
	grpc := NewGRPCServer()
	grpc.Start()

	StartQueue(context.Background(), 2)
	go func() {
		sig := <-sigs
		//bm.Stop()
		fmt.Println(sig, "exiting")
		StopQueue()
		done <- true
	}()
	<-done
}
