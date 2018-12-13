package main

import "github.com/godcong/go-ffmpeg/service"

//go:generate swagger generate spec
func main() {
	service.RunMain()
}
