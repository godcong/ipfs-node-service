package service

import (
	"context"
	"fmt"
	"github.com/godcong/ipfs-media-service/config"
	"github.com/godcong/ipfs-media-service/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"syscall"
	"time"
)

// GRPCServer ...
type GRPCServer struct {
	config *config.Configure
	server *grpc.Server
	Type   string
	Port   string
	Path   string
}

// RemoteDownload ...
func (s *GRPCServer) RemoteDownload(ctx context.Context, p *proto.RemoteDownloadRequest) (*proto.NodeReply, error) {
	log.Printf("Received: %v", p.String())
	stream := NewStreamerWithConfig(s.config, p.ID)
	//stream.Dir, stream.FileName = filepath.Split(key)
	stream.ObjectKey = p.ObjectKey
	stream.SetEncrypt(false)
	stream.Callback = s.config.Callback.Type
	//stream.SetURI("")
	//stream.FileDest = config.Media.Upload
	//stream.SetSrc(config.Media.Transfer)
	globalQueue.Set(stream.ID, StatusQueuing, 0)
	Push(stream)
	return Result(nil), nil
}

type grpcBack struct {
	config   *config.Configure
	BackType string
	BackAddr string
}

// NewManagerGRPC ...
func NewManagerGRPC(cfg *config.Configure) *GRPCClient {
	return &GRPCClient{
		config: cfg,
		Type:   config.DefaultString("unix", Type),
		Port:   config.DefaultString("", ":7781"),
		Addr:   config.DefaultString("", "/tmp/manager.sock"),
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
		conn, err = grpc.Dial(c.Addr, grpc.WithInsecure())
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

// Callback ...
func (b *grpcBack) Callback(r *QueueResult) error {
	grpc := NewManagerGRPC(config.Config())
	client := ManagerClient(grpc)
	timeout, _ := context.WithTimeout(context.Background(), time.Second*5)

	reply, err := client.CensorBack(timeout, &proto.ManagerCensorRequest{
		ID:     r.ID,
		Detail: r.JSON(),
	})
	if err != nil {
		log.Error(err)
	}
	log.Println(reply)
	return err
}

// NewGRPCBack ...
func NewGRPCBack(cfg *config.Configure) StreamerCallback {
	return &grpcBack{
		config:   cfg,
		BackType: config.DefaultString(cfg.Callback.BackType, "tcp"),
		BackAddr: config.DefaultString(cfg.Callback.BackAddr, "localhost:7781"),
	}
}

// Status ...
func (s *GRPCServer) Status(ctx context.Context, p *proto.StatusRequest) (*proto.NodeReply, error) {
	log.Printf("Received: %v", p.String())
	return Result(nil), nil
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
	s.server = grpc.NewServer()
	var lis net.Listener
	var port string
	var err error
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

		proto.RegisterNodeServiceServer(s.server, s)
		// Register reflection service on gRPC server.
		reflection.Register(s.server)
		log.Printf("Listening and serving TCP on %s\n", port)
		if err := s.server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

}

// Stop ...
func (s *GRPCServer) Stop() {
	s.server.Stop()
}
