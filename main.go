//go:generate apidoc -i ./service
//go:generate statik -f -src=./doc
//go:generate protoc --go_out=plugins=grpc:./proto node.proto
package main

import (
	"flag"
	"github.com/godcong/node-service/config"
	"github.com/godcong/node-service/service"
	_ "github.com/godcong/node-service/statik"
	"io"
	"log"
	"os"
)

var configPath = flag.String("path", "config.toml", "load config file from path")

func main() {
	flag.Parse()

	file, err := os.OpenFile("node.log", os.O_SYNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	log.SetOutput(io.MultiWriter(file, os.Stdout))
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err = config.Initialize(*configPath)
	if err != nil {
		panic(err)
	}

	service.RunMain()
}
