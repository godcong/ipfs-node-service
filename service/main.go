package service

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// RunMain 主线程
func RunMain() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := Initialize(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//new rest
	serv := NewRestServer(":8060")
	_ = Router(serv.Engine)

	//start
	serv.Start()
	StartQueue(2)
	go func() {
		sig := <-sigs
		//bm.Stop()
		fmt.Println(sig, "exiting")
		StopQueue()
		done <- true
	}()
	<-done
}
