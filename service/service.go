package service

import (
	"github.com/godcong/ipfs-node-service/config"
	"github.com/godcong/ipfs-node-service/oss"
	log "github.com/sirupsen/logrus"
)

// service ...
type service struct {
	rest  *RestServer
	grpc  *GRPCServer
	queue *QueueServer
}

var server *service

// Start 主线程
func Start() {
	cfg := config.Config()

	server = &service{
		grpc:  NewGRPCServer(cfg),
		rest:  NewRestServer(cfg),
		queue: NewQueueServer(cfg),
	}

	log.Println("run main")
	oss.InitOSS(cfg)

	server.rest.Start()
	server.grpc.Start()

	//queue start
	server.queue.Processes = 5
	server.queue.Start()

}

// Stop ...
func Stop() {
	server.rest.Stop()
	server.grpc.Stop()
	server.queue.Stop()
}

// NewBack ...
func NewBack() StreamerCallback {
	cfg := config.Config()
	if cfg != nil && cfg.Node.RequestType == "rest" {
		return NewRestBack(cfg)
	}
	return NewGRPCBack(cfg)
}
