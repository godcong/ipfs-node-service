//go:generate apidoc -i ./service
//go:generate statik -f -src=./doc
package main

import (
	"github.com/godcong/go-ffmpeg/service"
	"log"
	"os"
)
import _ "github.com/godcong/go-ffmpeg/statik"

func main() {
	file, err := os.OpenFile("node.log", os.O_SYNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	service.RunMain()
}
