package service

import (
	"context"
	"fmt"
	"github.com/godcong/ipfs-media-service/config"
	"github.com/godcong/ipfs-media-service/proto"
	"github.com/micro/go-micro"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"syscall"
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
	stream.Callback = s.config.ManagerConfig.CallType
	//stream.SetURI("")
	//stream.FileDest = config.Media.Upload
	//stream.SetSrc(config.Media.Transfer)
	globalQueue.Set(stream.ID, StatusQueuing, 0)
	Push(stream)
	rep = Result(&proto.NodeReplyDetail{
		ID: stream.ID,
	})
	return nil
}

func (s *GRPCServer) Status(ctx context.Context, req *proto.StatusRequest, rep *proto.NodeReply) error {
	log.Printf("Received: %v", req.String())
	rep = Result(nil)
	return nil
}

// NewManagerGRPC ...
func NewManagerGRPC(cfg *config.Configure) *GRPCClient {

	return &GRPCClient{
		config: cfg,
		Type:   config.DefaultString("tcp", Type),
		Port:   config.DefaultString("", ":7781"),
		Addr:   config.DefaultString("", "localhost"),
	}
}

// GRPCClient ...
type GRPCClient struct {
	config *config.Configure
	Type   string
	Port   string
	Addr   string
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
func Result(detail *proto.NodeReplyDetail) *proto.NodeReply {
	return &proto.NodeReply{
		Code:    0,
		Message: "success",
		Detail:  detail,
	}
}

// NewGRPCServer ...
func NewGRPCServer(cfg *config.Configure) *GRPCServer {
	return &GRPCServer{
		config: cfg,
		Type:   config.DefaultString(cfg.GRPC.Type, Type),
		Port:   config.DefaultString(cfg.GRPC.Port, ":7788"),
		Path:   config.DefaultString(cfg.GRPC.Path, "/tmp/node.sock"),
	}
}

// Start ...
func (s *GRPCServer) Start() {
	if !s.config.GRPC.Enable {
		return
	}
	var lis net.Listener
	var port string
	var err error
	//reg := consul.NewRegistry()

	s.service = micro.NewService(
		micro.Name("go.micro.grpc.node"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.Version("latest"),
	)
	s.service.Init()
	go func() {
		if s.Type == "unix" {
			_ = syscall.Unlink(s.Path)
			lis, err = net.Listen(s.Type, s.Path)
			port = s.Path
		} else {
			lis, err = net.Listen("tcp", s.Port)
			port = s.Port
		}

		if err != nil {
			panic(fmt.Sprintf("failed to listen: %v", err))
		}

		_ = proto.RegisterNodeServiceHandler(s.service.Server(), s)

		log.Printf("Listening and serving TCP on %s\n", port)
		if err := s.service.Run(); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

}

// Stop ...
func (s *GRPCServer) Stop() {
	//s.server.Stop()
}
