//go:generate apidoc -i ./service
//go:generate statik -f -src=./doc
//go:generate protoc --go_out=plugins=grpc:./proto --micro_out=./proto node.proto
package main

import (
	"flag"
	"github.com/godcong/go-trait"
	"github.com/godcong/ipfs-node-service/config"
	"github.com/godcong/ipfs-node-service/service"
	_ "github.com/godcong/ipfs-node-service/statik"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var configPath = flag.String("config", "config.toml", "load config file from path")
var elk = flag.Bool("elk", false, "set log to elk")
var logPath = flag.String("log", "logs/manager.log", "set the default log path")

func main() {
	flag.Parse()

	if *elk {
		trait.InitElasticLog("ipfs-node-service", nil)
	} else {
		trait.InitRotateLog(*logPath, nil)
	}

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
		log.Info(sig, "exiting")
		service.Stop()
		done <- true
	}()
	<-done

}
