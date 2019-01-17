package service

import (
	"context"
	"fmt"
	"github.com/godcong/node-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

// GRPCServer ...
type GRPCServer struct {
	Port   string
	server *grpc.Server
}

// RemoteDownload ...
func (s *GRPCServer) RemoteDownload(ctx context.Context, p *proto.RemoteDownloadRequest) (*proto.ServiceReply, error) {
	log.Printf("Received: %v", p.String())
	return Result(nil), nil
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
	port := config.GRPC.Port
	if port == "" {
		port = ":7782"
	}
	return &GRPCServer{
		Port: port,
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

		if config.GRPC.Type == "sock" {
			_ = os.Remove("/tmp/node.sock")
			lis, err = net.Listen("unix", "/tmp/node.sock")
			port = lis.Addr().String()
		} else {
			lis, err = net.Listen("tcp", config.GRPC.Port)
			port = config.GRPC.Port
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
