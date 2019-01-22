// Code generated by protoc-gen-go. DO NOT EDIT.
// source: node.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// NodeRemoteType ...
type NodeRemoteType int32

// NodeRemoteType_Basic ...
const (
	NodeRemoteType_Basic NodeRemoteType = 0
	NodeRemoteType_Retry NodeRemoteType = 1
	NodeRemoteType_Force NodeRemoteType = 2
)

// NodeRemoteType_name ...
var NodeRemoteType_name = map[int32]string{
	0: "Basic",
	1: "Retry",
	2: "Force",
}

// NodeRemoteType_value ...
var NodeRemoteType_value = map[string]int32{
	"Basic": 0,
	"Retry": 1,
	"Force": 2,
}

// String ...
func (x NodeRemoteType) String() string {
	return proto.EnumName(NodeRemoteType_name, int32(x))
}

// EnumDescriptor ...
func (NodeRemoteType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{0}
}

// NodeBackType ...
type NodeBackType int32

// NodeBackType_HTTP ...
const (
	NodeBackType_HTTP NodeBackType = 0
	NodeBackType_GRPC NodeBackType = 1
)

// NodeBackType_name ...
var NodeBackType_name = map[int32]string{
	0: "HTTP",
	1: "GRPC",
}

// NodeBackType_value ...
var NodeBackType_value = map[string]int32{
	"HTTP": 0,
	"GRPC": 1,
}

// String ...
func (x NodeBackType) String() string {
	return proto.EnumName(NodeBackType_name, int32(x))
}

// EnumDescriptor ...
func (NodeBackType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{1}
}

// StatusRequest ...
type StatusRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

// Reset ...
func (m *StatusRequest) Reset() { *m = StatusRequest{} }

// String ...
func (m *StatusRequest) String() string { return proto.CompactTextString(m) }

// ProtoMessage ...
func (*StatusRequest) ProtoMessage() {}

// Descriptor ...
func (*StatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{0}
}

// XXX_Unmarshal ...
func (m *StatusRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusRequest.Unmarshal(m, b)
}

// XXX_Marshal ...
func (m *StatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusRequest.Marshal(b, m, deterministic)
}

// XXX_Merge ...
func (m *StatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusRequest.Merge(m, src)
}

// XXX_Size ...
func (m *StatusRequest) XXX_Size() int {
	return xxx_messageInfo_StatusRequest.Size(m)
}

// XXX_DiscardUnknown ...
func (m *StatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StatusRequest proto.InternalMessageInfo

// GetId ...
func (m *StatusRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

// RemoteDownloadRequest ...
type RemoteDownloadRequest struct {
	ObjectKey            string   `protobuf:"bytes,1,opt,name=objectKey,proto3" json:"objectKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

// Reset ...
func (m *RemoteDownloadRequest) Reset() { *m = RemoteDownloadRequest{} }

// String ...
func (m *RemoteDownloadRequest) String() string { return proto.CompactTextString(m) }

// ProtoMessage ...
func (*RemoteDownloadRequest) ProtoMessage() {}

// Descriptor ...
func (*RemoteDownloadRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{1}
}

// XXX_Unmarshal ...
func (m *RemoteDownloadRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoteDownloadRequest.Unmarshal(m, b)
}

// XXX_Marshal ...
func (m *RemoteDownloadRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoteDownloadRequest.Marshal(b, m, deterministic)
}

// XXX_Merge ...
func (m *RemoteDownloadRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoteDownloadRequest.Merge(m, src)
}

// XXX_Size ...
func (m *RemoteDownloadRequest) XXX_Size() int {
	return xxx_messageInfo_RemoteDownloadRequest.Size(m)
}

// XXX_DiscardUnknown ...
func (m *RemoteDownloadRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoteDownloadRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RemoteDownloadRequest proto.InternalMessageInfo

// GetObjectKey ...
func (m *RemoteDownloadRequest) GetObjectKey() string {
	if m != nil {
		return m.ObjectKey
	}
	return ""
}

// NodeReply ...
type NodeReply struct {
	Code                 int32            `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message              string           `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Detail               *NodeReplyDetail `protobuf:"bytes,3,opt,name=detail,proto3" json:"detail,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

// Reset ...
func (m *NodeReply) Reset() { *m = NodeReply{} }

// String ...
func (m *NodeReply) String() string { return proto.CompactTextString(m) }

// ProtoMessage ...
func (*NodeReply) ProtoMessage() {}

// Descriptor ...
func (*NodeReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{2}
}

// XXX_Unmarshal ...
func (m *NodeReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeReply.Unmarshal(m, b)
}

// XXX_Marshal ...
func (m *NodeReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeReply.Marshal(b, m, deterministic)
}

// XXX_Merge ...
func (m *NodeReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeReply.Merge(m, src)
}

// XXX_Size ...
func (m *NodeReply) XXX_Size() int {
	return xxx_messageInfo_NodeReply.Size(m)
}

// XXX_DiscardUnknown ...
func (m *NodeReply) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeReply.DiscardUnknown(m)
}

var xxx_messageInfo_NodeReply proto.InternalMessageInfo

// GetCode ...
func (m *NodeReply) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

// GetMessage ...
func (m *NodeReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

// GetDetail ...
func (m *NodeReply) GetDetail() *NodeReplyDetail {
	if m != nil {
		return m.Detail
	}
	return nil
}

// NodeReplyDetail ...
type NodeReplyDetail struct {
	ID                   string   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Json                 string   `protobuf:"bytes,2,opt,name=json,proto3" json:"json,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

// Reset ...
func (m *NodeReplyDetail) Reset() { *m = NodeReplyDetail{} }

// String ...
func (m *NodeReplyDetail) String() string { return proto.CompactTextString(m) }

// ProtoMessage ...
func (*NodeReplyDetail) ProtoMessage() {}

// Descriptor ...
func (*NodeReplyDetail) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{3}
}

// XXX_Unmarshal ...
func (m *NodeReplyDetail) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeReplyDetail.Unmarshal(m, b)
}

// XXX_Marshal ...
func (m *NodeReplyDetail) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeReplyDetail.Marshal(b, m, deterministic)
}

// XXX_Merge ...
func (m *NodeReplyDetail) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeReplyDetail.Merge(m, src)
}

// XXX_Size ...
func (m *NodeReplyDetail) XXX_Size() int {
	return xxx_messageInfo_NodeReplyDetail.Size(m)
}

// XXX_DiscardUnknown ...
func (m *NodeReplyDetail) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeReplyDetail.DiscardUnknown(m)
}

var xxx_messageInfo_NodeReplyDetail proto.InternalMessageInfo

// GetID ...
func (m *NodeReplyDetail) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

// GetJson ...
func (m *NodeReplyDetail) GetJson() string {
	if m != nil {
		return m.Json
	}
	return ""
}

func init() {
	proto.RegisterEnum("proto.NodeRemoteType", NodeRemoteType_name, NodeRemoteType_value)
	proto.RegisterEnum("proto.NodeBackType", NodeBackType_name, NodeBackType_value)
	proto.RegisterType((*StatusRequest)(nil), "proto.StatusRequest")
	proto.RegisterType((*RemoteDownloadRequest)(nil), "proto.RemoteDownloadRequest")
	proto.RegisterType((*NodeReply)(nil), "proto.NodeReply")
	proto.RegisterType((*NodeReplyDetail)(nil), "proto.NodeReplyDetail")
}

func init() { proto.RegisterFile("node.proto", fileDescriptor_0c843d59d2d938e7) }

var fileDescriptor_0c843d59d2d938e7 = []byte{
	// 340 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x51, 0xd1, 0x4e, 0xea, 0x40,
	0x10, 0xa5, 0xbd, 0xc0, 0xbd, 0x1d, 0xee, 0xe5, 0x36, 0x1b, 0x35, 0x8d, 0x21, 0x91, 0xf4, 0x89,
	0xf0, 0xd0, 0x44, 0x0c, 0x3f, 0x50, 0x1b, 0x95, 0x98, 0x98, 0x66, 0xe1, 0x07, 0xca, 0xee, 0xa4,
	0x16, 0xa1, 0x83, 0xed, 0xa2, 0xe9, 0x07, 0xf8, 0xdf, 0x66, 0xb7, 0x45, 0x85, 0xf0, 0xb4, 0x67,
	0xe7, 0x9c, 0xb3, 0x3b, 0x67, 0x06, 0x20, 0x27, 0x89, 0xc1, 0xb6, 0x20, 0x45, 0xac, 0x63, 0x0e,
	0xff, 0x0a, 0xfe, 0xcd, 0x55, 0xa2, 0x76, 0x25, 0xc7, 0xd7, 0x1d, 0x96, 0x8a, 0xf5, 0xc1, 0xce,
	0xa4, 0x67, 0x0d, 0xad, 0x91, 0xc3, 0xed, 0x4c, 0xfa, 0x53, 0x38, 0xe7, 0xb8, 0x21, 0x85, 0x11,
	0xbd, 0xe7, 0x6b, 0x4a, 0xe4, 0x5e, 0x38, 0x00, 0x87, 0x96, 0x2b, 0x14, 0xea, 0x11, 0xab, 0x46,
	0xff, 0x5d, 0xf0, 0x33, 0x70, 0x9e, 0x48, 0x22, 0xc7, 0xed, 0xba, 0x62, 0x0c, 0xda, 0x82, 0x24,
	0x1a, 0x55, 0x87, 0x1b, 0xcc, 0x3c, 0xf8, 0xbd, 0xc1, 0xb2, 0x4c, 0x52, 0xf4, 0x6c, 0x63, 0xde,
	0x5f, 0x59, 0x00, 0x5d, 0x89, 0x2a, 0xc9, 0xd6, 0xde, 0xaf, 0xa1, 0x35, 0xea, 0x4d, 0x2e, 0xea,
	0x8e, 0x83, 0xaf, 0xf7, 0x22, 0xc3, 0xf2, 0x46, 0xe5, 0x4f, 0xe1, 0xff, 0x11, 0xa5, 0x43, 0xcc,
	0xa2, 0x7d, 0x88, 0x59, 0xa4, 0x1b, 0x58, 0x95, 0x94, 0x37, 0x3f, 0x19, 0x3c, 0xbe, 0x86, 0x7e,
	0x6d, 0xd3, 0xe1, 0x16, 0xd5, 0x16, 0x99, 0x03, 0x9d, 0x30, 0x29, 0x33, 0xe1, 0xb6, 0x34, 0xe4,
	0xa8, 0x8a, 0xca, 0xb5, 0x34, 0xbc, 0xa3, 0x42, 0xa0, 0x6b, 0x8f, 0x7d, 0xf8, 0xab, 0x2d, 0x61,
	0x22, 0x5e, 0x8c, 0xe1, 0x0f, 0xb4, 0x1f, 0x16, 0x8b, 0xd8, 0x6d, 0x69, 0x74, 0xcf, 0xe3, 0x5b,
	0xd7, 0x9a, 0x7c, 0x58, 0xd0, 0xd3, 0xa2, 0x39, 0x16, 0x6f, 0x99, 0x40, 0x16, 0x42, 0xff, 0x70,
	0x7e, 0x6c, 0xd0, 0xe4, 0x39, 0x39, 0xd6, 0x4b, 0xf7, 0x38, 0xad, 0xdf, 0x62, 0x13, 0xe8, 0xd6,
	0x4b, 0x62, 0x67, 0x0d, 0x7b, 0xb0, 0xb3, 0x53, 0x9e, 0x30, 0x00, 0x4f, 0xd0, 0x26, 0x48, 0x33,
	0xf5, 0xbc, 0x5b, 0x06, 0x29, 0x49, 0x41, 0x79, 0x5a, 0xeb, 0x42, 0xf7, 0x47, 0x83, 0xb1, 0xae,
	0xc4, 0xd6, 0xb2, 0x6b, 0xa8, 0x9b, 0xcf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x9d, 0xfd, 0x22, 0x4d,
	0x24, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// NodeServiceClient is the client API for NodeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NodeServiceClient interface {
	RemoteDownload(ctx context.Context, in *RemoteDownloadRequest, opts ...grpc.CallOption) (*NodeReply, error)
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*NodeReply, error)
}

type nodeServiceClient struct {
	cc *grpc.ClientConn
}

// NewNodeServiceClient ...
func NewNodeServiceClient(cc *grpc.ClientConn) NodeServiceClient {
	return &nodeServiceClient{cc}
}

// RemoteDownload ...
func (c *nodeServiceClient) RemoteDownload(ctx context.Context, in *RemoteDownloadRequest, opts ...grpc.CallOption) (*NodeReply, error) {
	out := new(NodeReply)
	err := c.cc.Invoke(ctx, "/proto.NodeService/RemoteDownload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Status ...
func (c *nodeServiceClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*NodeReply, error) {
	out := new(NodeReply)
	err := c.cc.Invoke(ctx, "/proto.NodeService/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeServiceServer is the server API for NodeService service.
type NodeServiceServer interface {
	RemoteDownload(context.Context, *RemoteDownloadRequest) (*NodeReply, error)
	Status(context.Context, *StatusRequest) (*NodeReply, error)
}

// RegisterNodeServiceServer ...
func RegisterNodeServiceServer(s *grpc.Server, srv NodeServiceServer) {
	s.RegisterService(&_NodeService_serviceDesc, srv)
}

func _NodeService_RemoteDownload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoteDownloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServiceServer).RemoteDownload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NodeService/RemoteDownload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServiceServer).RemoteDownload(ctx, req.(*RemoteDownloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NodeService_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServiceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NodeService/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServiceServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _NodeService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.NodeService",
	HandlerType: (*NodeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RemoteDownload",
			Handler:    _NodeService_RemoteDownload_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _NodeService_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "node.proto",
}
