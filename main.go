//go:generate apidoc -i ./service
//go:generate statik -f -src=./doc
//go:generate protoc --go_out=plugins=grpc:./proto node.proto
package main

import (
	"flag"
	"fmt"
	"github.com/godcong/go-trait"
	"github.com/godcong/ipfs-media-service/config"
	"github.com/godcong/ipfs-media-service/service"
	_ "github.com/godcong/ipfs-media-service/statik"
	"os"
	"os/signal"
	"syscall"
)

var configPath = flag.String("path", "config.toml", "load config file from path")

func main() {
	flag.Parse()
	trait.InitElasticLog("ipfs-node-service", nil)

	err := config.Initialize(*configPath)
	if err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//start
	service.Start()

	go func() {
		sig := <-sigs
		//bm.Stop()
		fmt.Println(sig, "exiting")
		service.Stop()
		done <- true
	}()
	<-done

}
