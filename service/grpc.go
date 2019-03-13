package service

import (
	"context"
	"github.com/godcong/ipfs-node-service/config"
	"github.com/godcong/ipfs-node-service/proto"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry/consul"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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
	//stream.FileDest = config.Media.Upload
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

// NewManagerGRPC ...
func NewManagerGRPC(cfg *config.Configure) *GRPCClient {
	reg := consul.NewRegistry()
	return &GRPCClient{
		config: cfg,
		service: micro.NewService(
			micro.Registry(reg),
		),
		Type: config.DefaultString("tcp", Type),
		Port: config.DefaultString("", ":7781"),
		Addr: config.DefaultString("", "localhost"),
	}
}

// GRPCClient ...
type GRPCClient struct {
	config  *config.Configure
	service micro.Service
	Type    string
	Port    string
	Addr    string
}

// Conn ...
func (c *GRPCClient) Conn() (*grpc.ClientConn, error) {

	var conn *grpc.ClientConn
	var err error

	if c.Type == "unix" {
		conn, err = grpc.Dial("passthrough:///unix://"+c.Addr, grpc.WithInsecure())
	} else {
		conn, err = grpc.Dial(c.Addr+c.Port, grpc.WithInsecure())
	}

	return conn, err
}

// ManagerClient ...
func ManagerClient(g *GRPCClient) proto.ManagerServiceClient {
	clientConn, err := g.Conn()
	if err != nil {
		log.Println(err)
		return nil
	}
	return proto.NewManagerServiceClient(clientConn)
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
		micro.Name("go.micro.grpc.node"),
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
