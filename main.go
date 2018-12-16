//go:generate apidoc -i ./service
//go:generate statik -f -src=./doc
package main

import "github.com/godcong/go-ffmpeg/service"
import _ "github.com/godcong/go-ffmpeg/statik"

func main() {
	service.RunMain()
}
