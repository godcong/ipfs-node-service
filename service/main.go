package service

import (
	"github.com/godcong/node-service/config"
	"github.com/godcong/node-service/oss"
	"log"
)

// Service ...
type Service struct {
	rest  *RestServer
	grpc  *GRPCServer
	queue *QueueServer
}

var server *Service

// Start 主线程
func Start() {
	cfg := config.Config()

	server = &Service{
		grpc:  NewGRPCServer(cfg),
		rest:  NewRestServer(cfg),
		queue: NewQueueServer(cfg),
	}

	log.Println("run main")

	oss.InitOSS(cfg)

	//rest start
	_ = Router(server.rest.Engine)
	server.rest.Start()

	//grpc start
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
	if cfg != nil && cfg.Callback.Type == "grpc" {
		return NewGRPCBack(cfg)
	}
	return NewRestBack(cfg)
}
