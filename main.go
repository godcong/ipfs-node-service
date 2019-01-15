//go:generate apidoc -i ./service
//go:generate statik -f -src=./doc
//go:generate protoc --go_out=plugins=grpc:./proto node.proto
package main

import (
	"flag"
	"github.com/godcong/node-service/service"
	_ "github.com/godcong/node-service/statik"
	"log"
	"os"
	"runtime/pprof"
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}
	file, err := os.OpenFile("node.log", os.O_SYNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	service.RunMain()
}
