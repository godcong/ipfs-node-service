package service

import (
	"github.com/godcong/go-ffmpeg/ipfs"
	"github.com/ipfs/go-ipfs-api"
)

var api = InitIPFS("localhost", "5001")

// InitIPFS ...
func InitIPFS(url, port string) ipfs.API {
	return ipfs.NewConfig(url + ":" + port).VersionAPI("v0")
}

// Shell ...
func Shell() *shell.Shell {
	return shell.NewShell("localhost:5001")
}
