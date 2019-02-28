package service

import (
	"github.com/godcong/ipfs-media-service/ipfs"
	log "github.com/sirupsen/logrus"
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

func commitToIPNS(id, source string) (map[string]interface{}, error) {
	api = DefaultIPFS()
	fs, err := api.AddDir(source)
	if err != nil {
		return nil, err
	}
	log.Println(fs)

	if id != "" || err != nil {
		m, _ := api.Key().Gen(id, "rsa", 2048)
		if err != nil {
			//ignore error
		}
		log.Println(m, err)
	}
	ns, err := api.Name().PublishWithKey("/ipfs/"+fs["Hash"], id)
	if err != nil {
		return nil, err
	}
	log.Println(ns)

	return map[string]interface{}{
		"ID":      id,
		"fs_info": fs,
		"ns_info": ns,
	}, nil
}
