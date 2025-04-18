// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// MessagesClient is the client API for Messages service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessagesClient interface {
	GetServers(ctx context.Context, in *GetServersRequest, opts ...grpc.CallOption) (*GetServersResponse, error)
	SaveMessage(ctx context.Context, in *SaveMessageRequest, opts ...grpc.CallOption) (*SaveMessageResponse, error)
	DeleteMessage(ctx context.Context, in *DeleteMessageRequest, opts ...grpc.CallOption) (*DeleteMessageResponse, error)
	UpdateMessage(ctx context.Context, in *UpdateMessageRequest, opts ...grpc.CallOption) (*UpdateMessageResponse, error)
	ReadOneMessage(ctx context.Context, in *ReadOneMessageRequest, opts ...grpc.CallOption) (*Message, error)
	ReadAllMessages(ctx context.Context, in *ReadMessagesRequest, opts ...grpc.CallOption) (*ReadMessagesResponse, error)
	ReadThreadMessages(ctx context.Context, in *ReadThreadMessagesRequest, opts ...grpc.CallOption) (*ReadThreadMessagesResponse, error)
}

type messagesClient struct {
	cc grpc.ClientConnInterface
}

func NewMessagesClient(cc grpc.ClientConnInterface) MessagesClient {
	return &messagesClient{cc}
}

func (c *messagesClient) GetServers(ctx context.Context, in *GetServersRequest, opts ...grpc.CallOption) (*GetServersResponse, error) {
	out := new(GetServersResponse)
	err := c.cc.Invoke(ctx, "/messages.v1.Messages/GetServers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) SaveMessage(ctx context.Context, in *SaveMessageRequest, opts ...grpc.CallOption) (*SaveMessageResponse, error) {
	out := new(SaveMessageResponse)
	err := c.cc.Invoke(ctx, "/messages.v1.Messages/SaveMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) DeleteMessage(ctx context.Context, in *DeleteMessageRequest, opts ...grpc.CallOption) (*DeleteMessageResponse, error) {
	out := new(DeleteMessageResponse)
	err := c.cc.Invoke(ctx, "/messages.v1.Messages/DeleteMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) UpdateMessage(ctx context.Context, in *UpdateMessageRequest, opts ...grpc.CallOption) (*UpdateMessageResponse, error) {
	out := new(UpdateMessageResponse)
	err := c.cc.Invoke(ctx, "/messages.v1.Messages/UpdateMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) ReadOneMessage(ctx context.Context, in *ReadOneMessageRequest, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/messages.v1.Messages/ReadOneMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) ReadAllMessages(ctx context.Context, in *ReadMessagesRequest, opts ...grpc.CallOption) (*ReadMessagesResponse, error) {
	out := new(ReadMessagesResponse)
	err := c.cc.Invoke(ctx, "/messages.v1.Messages/ReadAllMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) ReadThreadMessages(ctx context.Context, in *ReadThreadMessagesRequest, opts ...grpc.CallOption) (*ReadThreadMessagesResponse, error) {
	out := new(ReadThreadMessagesResponse)
	err := c.cc.Invoke(ctx, "/messages.v1.Messages/ReadThreadMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessagesServer is the server API for Messages service.
// All implementations must embed UnimplementedMessagesServer
// for forward compatibility
type MessagesServer interface {
	GetServers(context.Context, *GetServersRequest) (*GetServersResponse, error)
	SaveMessage(context.Context, *SaveMessageRequest) (*SaveMessageResponse, error)
	DeleteMessage(context.Context, *DeleteMessageRequest) (*DeleteMessageResponse, error)
	UpdateMessage(context.Context, *UpdateMessageRequest) (*UpdateMessageResponse, error)
	ReadOneMessage(context.Context, *ReadOneMessageRequest) (*Message, error)
	ReadAllMessages(context.Context, *ReadMessagesRequest) (*ReadMessagesResponse, error)
	ReadThreadMessages(context.Context, *ReadThreadMessagesRequest) (*ReadThreadMessagesResponse, error)
	mustEmbedUnimplementedMessagesServer()
}

// UnimplementedMessagesServer must be embedded to have forward compatible implementations.
type UnimplementedMessagesServer struct {
}

func (UnimplementedMessagesServer) GetServers(context.Context, *GetServersRequest) (*GetServersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetServers not implemented")
}
func (UnimplementedMessagesServer) SaveMessage(context.Context, *SaveMessageRequest) (*SaveMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveMessage not implemented")
}
func (UnimplementedMessagesServer) DeleteMessage(context.Context, *DeleteMessageRequest) (*DeleteMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMessage not implemented")
}
func (UnimplementedMessagesServer) UpdateMessage(context.Context, *UpdateMessageRequest) (*UpdateMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMessage not implemented")
}
func (UnimplementedMessagesServer) ReadOneMessage(context.Context, *ReadOneMessageRequest) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadOneMessage not implemented")
}
func (UnimplementedMessagesServer) ReadAllMessages(context.Context, *ReadMessagesRequest) (*ReadMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadAllMessages not implemented")
}
func (UnimplementedMessagesServer) ReadThreadMessages(context.Context, *ReadThreadMessagesRequest) (*ReadThreadMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadThreadMessages not implemented")
}
func (UnimplementedMessagesServer) mustEmbedUnimplementedMessagesServer() {}

// UnsafeMessagesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessagesServer will
// result in compilation errors.
type UnsafeMessagesServer interface {
	mustEmbedUnimplementedMessagesServer()
}

func RegisterMessagesServer(s *grpc.Server, srv MessagesServer) {
	s.RegisterService(&_Messages_serviceDesc, srv)
}

func _Messages_GetServers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetServersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagesServer).GetServers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.v1.Messages/GetServers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagesServer).GetServers(ctx, req.(*GetServersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messages_SaveMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagesServer).SaveMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.v1.Messages/SaveMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagesServer).SaveMessage(ctx, req.(*SaveMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messages_DeleteMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagesServer).DeleteMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.v1.Messages/DeleteMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagesServer).DeleteMessage(ctx, req.(*DeleteMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messages_UpdateMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagesServer).UpdateMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.v1.Messages/UpdateMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagesServer).UpdateMessage(ctx, req.(*UpdateMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messages_ReadOneMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadOneMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagesServer).ReadOneMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.v1.Messages/ReadOneMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagesServer).ReadOneMessage(ctx, req.(*ReadOneMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messages_ReadAllMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagesServer).ReadAllMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.v1.Messages/ReadAllMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagesServer).ReadAllMessages(ctx, req.(*ReadMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messages_ReadThreadMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadThreadMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagesServer).ReadThreadMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.v1.Messages/ReadThreadMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagesServer).ReadThreadMessages(ctx, req.(*ReadThreadMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Messages_serviceDesc = grpc.ServiceDesc{
	ServiceName: "messages.v1.Messages",
	HandlerType: (*MessagesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetServers",
			Handler:    _Messages_GetServers_Handler,
		},
		{
			MethodName: "SaveMessage",
			Handler:    _Messages_SaveMessage_Handler,
		},
		{
			MethodName: "DeleteMessage",
			Handler:    _Messages_DeleteMessage_Handler,
		},
		{
			MethodName: "UpdateMessage",
			Handler:    _Messages_UpdateMessage_Handler,
		},
		{
			MethodName: "ReadOneMessage",
			Handler:    _Messages_ReadOneMessage_Handler,
		},
		{
			MethodName: "ReadAllMessages",
			Handler:    _Messages_ReadAllMessages_Handler,
		},
		{
			MethodName: "ReadThreadMessages",
			Handler:    _Messages_ReadThreadMessages_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/messages.proto",
}
