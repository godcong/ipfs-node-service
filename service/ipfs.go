package service

import "github.com/godcong/go-ffmpeg/ipfs"

var api = InitIPFS("localhost", "5001")

// InitIPFS ...
func InitIPFS(url, port string) ipfs.API {
	return ipfs.NewConfig(url + ":" + port).VersionAPI("v0")
}
