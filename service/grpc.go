package service

import (
	"context"
	"fmt"
	"github.com/godcong/node-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"syscall"
	"time"
)

// GRPCServer ...
type GRPCServer struct {
	Type   string
	Port   string
	Path   string
	server *grpc.Server
}

// RemoteDownload ...
func (s *GRPCServer) RemoteDownload(ctx context.Context, p *proto.RemoteDownloadRequest) (*proto.ServiceReply, error) {
	log.Printf("Received: %v", p.String())
	stream := NewStreamerWithConfig(Config(), p.ObjectKey)
	//stream.Dir, stream.FileName = filepath.Split(key)
	stream.ObjectKey = p.ObjectKey
	stream.SetEncrypt(false)
	stream.StreamerCallback = NewBack()
	//stream.SetURI("")
	//stream.FileDest = config.Media.Upload
	//stream.SetSrc(config.Media.Transfer)
	queue.Set(stream.ID, StatusQueuing, 0)
	Push(stream)
	return Result(nil), nil
}

type grpcBack struct {
	BackType string
	BackAddr string
}

// Callback ...
func (b *grpcBack) Callback(r *QueueResult) error {
	var conn *grpc.ClientConn
	var err error

	if b.BackType == "unix" {
		conn, err = grpc.Dial("passthrough:///unix://" + b.BackAddr)
		//conn, err = net.Dial(b.BackType, b.BackAddr)
	} else {
		conn, err = grpc.Dial(b.BackAddr)
		//conn, err = net.Dial("tcp", b.BackAddr)
	}
	if err != nil {
		return err
	}
	defer conn.Close()
	client := proto.NewManagerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	reply, err := client.Back(ctx, &proto.ManagerRequest{
		ObjectKey: r.ID,
		Detail:    r.JSON(),
	})

	if err != nil {
		return err
	}
	log.Printf("%+v", reply)
	return nil
}

// NewGRPCBack ...
func NewGRPCBack(cfg *Configure) StreamerCallback {
	return &grpcBack{
		BackType: DefaultString(cfg.GRPC.BackType, "tcp"),
		BackAddr: DefaultString(cfg.GRPC.BackAddr, "localhost:7783"),
	}
}

// Status ...
func (s *GRPCServer) Status(ctx context.Context, p *proto.StatusRequest) (*proto.ServiceReply, error) {
	log.Printf("Received: %v", p.String())
	return Result(nil), nil
}

// Result ...
func Result(detail *proto.ReplyDetail) *proto.ServiceReply {
	return &proto.ServiceReply{
		Code:                 0,
		Message:              "",
		Detail:               detail,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
}

// NewGRPCServer ...
func NewGRPCServer() *GRPCServer {
	return &GRPCServer{
		Type: DefaultString(config.GRPC.Type, Type),
		Port: DefaultString(config.GRPC.Port, ":7782"),
		Path: DefaultString(config.GRPC.Path, "/tmp/node.sock"),
	}
}

// Start ...
func (s *GRPCServer) Start() {
	if !config.GRPC.Enable {
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
