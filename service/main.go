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

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//new rest
	serv := NewRestServer(":8080")
	_ = Router(serv.Engine)

	//start
	serv.Start()
	StartQueue(5)
	go func() {
		sig := <-sigs
		//bm.Stop()
		fmt.Println(sig, "exiting")
		StopQueue()
		done <- true
	}()
	<-done
}
