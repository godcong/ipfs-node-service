package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// RunMain 主线程
func RunMain() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//new rest
	serv := NewRestServer(":8080")

	//start
	serv.Start()

	go func() {
		sig := <-sigs
		//bm.Stop()
		fmt.Println(sig, "exiting")
		done <- true
	}()
	<-done
}
