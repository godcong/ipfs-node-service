package service

import (
	"github.com/godcong/ipfs-node-service/ipfs"
	log "github.com/sirupsen/logrus"
	"path/filepath"
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
	log.Infof("commit info:%s,%s", id, source)
	api = DefaultIPFS()
	s, e := filepath.Abs(source)
	if e != nil {
		return nil, e
	}
	fs, e := api.AddDir(s)
	if e != nil {
		return nil, e
	}
	log.Println(fs)

	if id != "" || e != nil {
		m, _ := api.Key().Gen(id, "rsa", 2048)
		if e != nil {
			//ignore error
		}
		log.Println(m, e)
	}
	ns, e := api.Name().PublishWithKey("/ipfs/"+fs["Hash"], id)
	if e != nil {
		return nil, e
	}
	log.Println(ns)

	return map[string]interface{}{
		"ID":      id,
		"fs_info": fs,
		"ns_info": ns,
	}, nil
}
