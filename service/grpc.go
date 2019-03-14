package service

import (
	"context"
	"github.com/godcong/ipfs-node-service/config"
	"github.com/godcong/ipfs-node-service/proto"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry/consul"
	log "github.com/sirupsen/logrus"
	"time"
)

// GRPCServer ...
type GRPCServer struct {
	config  *config.Configure
	service micro.Service
	Type    string
	Port    string
	Path    string
}

func (s *GRPCServer) RemoteDownload(ctx context.Context, req *proto.RemoteDownloadRequest, rep *proto.NodeReply) error {
	log.Printf("Received: %v", req.String())
	stream := NewStreamerWithConfig(s.config, req.ID)
	//stream.Dir, stream.FileName = filepath.Split(key)
	stream.ObjectKey = req.ObjectKey
	stream.SetEncrypt(false)
	stream.Callback = s.config.Node.RequestType
	//stream.SetURI("")
	//stream.Transfer = config.Media.Upload
	//stream.SetSrc(config.Media.Transfer)
	globalQueue.Set(stream.ID, StatusQueuing, 0)
	Push(stream)
	*rep = Result(&proto.NodeReplyDetail{
		ID: stream.ID,
	})
	return nil
}

func (s *GRPCServer) Status(ctx context.Context, req *proto.StatusRequest, rep *proto.NodeReply) error {
	log.Printf("Received: %v", req.String())
	*rep = Result(nil)
	return nil
}

// GRPCClient ...
type GRPCClient struct {
	config  *config.Configure
	service micro.Service
}

// NewGRPCClient ...
func NewGRPCClient(cfg *config.Configure) *GRPCClient {
	reg := consul.NewRegistry()
	client := &GRPCClient{
		service: micro.NewService(
			micro.Registry(reg)),
		config: cfg,
	}
	client.service.Init()
	return client
}

// NodeClient ...
func NodeClient(g *GRPCClient) proto.NodeService {
	return proto.NewNodeService(g.config.Node.NodeName, g.service.Client())
}

// ManagerClient ...
func ManagerClient(g *GRPCClient) proto.ManagerService {
	return proto.NewManagerService(g.config.Node.ManagerName, g.service.Client())
}

// Result ...
func Result(detail *proto.NodeReplyDetail) proto.NodeReply {
	return proto.NodeReply{
		Code:    0,
		Message: "success",
		Detail:  detail,
	}
}

// NewGRPCServer ...
func NewGRPCServer(cfg *config.Configure) *GRPCServer {
	return &GRPCServer{
		config: cfg,
	}
}

// Start ...
func (s *GRPCServer) Start() {
	if !s.config.Node.EnableGRPC {
		return
	}

	reg := consul.NewRegistry()

	s.service = micro.NewService(
		micro.Name(s.config.Node.NodeName),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.Registry(reg),
	)
	s.service.Init()
	go func() {

		_ = proto.RegisterNodeServiceHandler(s.service.Server(), s)

		if err := s.service.Run(); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

}

// Stop ...
func (s *GRPCServer) Stop() {
	//s.server.Stop()
}
