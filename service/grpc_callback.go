package service

import (
	"context"
	"github.com/godcong/ipfs-media-service/config"
	"github.com/godcong/ipfs-media-service/proto"
	log "github.com/sirupsen/logrus"
	"time"
)

type grpcBack struct {
	config *config.Configure
}

// Callback ...
func (b *grpcBack) Callback(r *QueueResult) error {
	grpc := NewManagerGRPC(b.config)
	client := ManagerClient(grpc)
	timeout, _ := context.WithTimeout(context.Background(), time.Second*5)

	reply, err := client.NodeBack(timeout, &proto.ManagerNodeRequest{
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
