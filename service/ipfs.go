package service

import (
	"github.com/godcong/node-service/ipfs"
)

var api ipfs.API

// InitIPFS ...
func InitIPFS(url, port string) ipfs.API {
	api = ipfs.NewConfig(url + ":" + port).VersionAPI("v0")
	return api
}

// DefaultIPFS ...
func DefaultIPFS() ipfs.API {
	return InitIPFS("localhost", "5001")
}
