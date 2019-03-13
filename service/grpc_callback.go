package service

import (
	"context"
	"github.com/godcong/ipfs-node-service/config"
	"github.com/godcong/ipfs-node-service/proto"
	log "github.com/sirupsen/logrus"
	"time"
)

type grpcBack struct {
	config *config.Configure
}

// Callback ...
func (b *grpcBack) Callback(r *QueueResult) error {
	grpc := ManagerClient(NewGRPCClient(b.config))
	timeout, _ := context.WithTimeout(context.Background(), time.Second*5)
	reply, err := grpc.NodeBack(timeout, &proto.ManagerNodeRequest{
		ID:     r.ID,
		Detail: r.JSON(),
	})
	if err != nil {
		log.Error(err)
	}
	log.Printf("%+v,%+v", r, reply)
	return err
}

// NewGRPCBack ...
func NewGRPCBack(cfg *config.Configure) StreamerCallback {
	return &grpcBack{
		config: cfg,
	}
}
