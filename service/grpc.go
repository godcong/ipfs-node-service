package service

import (
	"context"
	"github.com/godcong/node-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// RemoteDownload ...
func (s *server) RemoteDownload(ctx context.Context, p *proto.RemoteDownloadRequest) (*proto.ServiceReply, error) {
	log.Printf("Received: %v", p.String())
	return Result(nil), nil
}

// Status ...
func (s *server) Status(ctx context.Context, p *proto.StatusRequest) (*proto.ServiceReply, error) {
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

// GRPCServerStart ...
func GRPCServerStart() {
	lis, err := net.Listen("tcp", ":7784")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	proto.RegisterNodeServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
