package service

import (
	"github.com/godcong/go-ffmpeg/ipfs"
)

var api ipfs.API

// InitIPFS ...
func InitIPFS(url, port string) ipfs.API {
	api = ipfs.NewConfig(url + ":" + port).VersionAPI("v0")
	return api
}

// Shell ...
//func Shell() *shell.Shell {
//	return shell.NewShell("localhost:5001")
//}
