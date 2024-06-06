// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.27.0
// source: chat.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	SendChat_SendChat_FullMethodName = "/SendChat/SendChat"
)

// SendChatClient is the client API for SendChat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SendChatClient interface {
	SendChat(ctx context.Context, in *ChatRequest, opts ...grpc.CallOption) (*ChatResponse, error)
}

type sendChatClient struct {
	cc grpc.ClientConnInterface
}

func NewSendChatClient(cc grpc.ClientConnInterface) SendChatClient {
	return &sendChatClient{cc}
}

func (c *sendChatClient) SendChat(ctx context.Context, in *ChatRequest, opts ...grpc.CallOption) (*ChatResponse, error) {
	out := new(ChatResponse)
	err := c.cc.Invoke(ctx, SendChat_SendChat_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SendChatServer is the server API for SendChat service.
// All implementations must embed UnimplementedSendChatServer
// for forward compatibility
type SendChatServer interface {
	SendChat(context.Context, *ChatRequest) (*ChatResponse, error)
	mustEmbedUnimplementedSendChatServer()
}

// UnimplementedSendChatServer must be embedded to have forward compatible implementations.
type UnimplementedSendChatServer struct {
}

func (UnimplementedSendChatServer) SendChat(context.Context, *ChatRequest) (*ChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendChat not implemented")
}
func (UnimplementedSendChatServer) mustEmbedUnimplementedSendChatServer() {}

// UnsafeSendChatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SendChatServer will
// result in compilation errors.
type UnsafeSendChatServer interface {
	mustEmbedUnimplementedSendChatServer()
}

func RegisterSendChatServer(s grpc.ServiceRegistrar, srv SendChatServer) {
	s.RegisterService(&SendChat_ServiceDesc, srv)
}

func _SendChat_SendChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SendChatServer).SendChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SendChat_SendChat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SendChatServer).SendChat(ctx, req.(*ChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SendChat_ServiceDesc is the grpc.ServiceDesc for SendChat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SendChat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SendChat",
	HandlerType: (*SendChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendChat",
			Handler:    _SendChat_SendChat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chat.proto",
}
