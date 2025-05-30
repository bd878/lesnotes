// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.19.4
// source: protos/messages.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Messages_GetServers_FullMethodName            = "/messages.v1.Messages/GetServers"
	Messages_SaveMessage_FullMethodName           = "/messages.v1.Messages/SaveMessage"
	Messages_DeleteMessage_FullMethodName         = "/messages.v1.Messages/DeleteMessage"
	Messages_DeleteMessages_FullMethodName        = "/messages.v1.Messages/DeleteMessages"
	Messages_DeleteAllUserMessages_FullMethodName = "/messages.v1.Messages/DeleteAllUserMessages"
	Messages_PublishMessages_FullMethodName       = "/messages.v1.Messages/PublishMessages"
	Messages_PrivateMessages_FullMethodName       = "/messages.v1.Messages/PrivateMessages"
	Messages_UpdateMessage_FullMethodName         = "/messages.v1.Messages/UpdateMessage"
	Messages_ReadOneMessage_FullMethodName        = "/messages.v1.Messages/ReadOneMessage"
	Messages_ReadAllMessages_FullMethodName       = "/messages.v1.Messages/ReadAllMessages"
	Messages_ReadThreadMessages_FullMethodName    = "/messages.v1.Messages/ReadThreadMessages"
)

// MessagesClient is the client API for Messages service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessagesClient interface {
	GetServers(ctx context.Context, in *GetServersRequest, opts ...grpc.CallOption) (*GetServersResponse, error)
	SaveMessage(ctx context.Context, in *SaveMessageRequest, opts ...grpc.CallOption) (*SaveMessageResponse, error)
	DeleteMessage(ctx context.Context, in *DeleteMessageRequest, opts ...grpc.CallOption) (*DeleteMessageResponse, error)
	DeleteMessages(ctx context.Context, in *DeleteMessagesRequest, opts ...grpc.CallOption) (*DeleteMessagesResponse, error)
	DeleteAllUserMessages(ctx context.Context, in *DeleteAllUserMessagesRequest, opts ...grpc.CallOption) (*DeleteAllUserMessagesResponse, error)
	PublishMessages(ctx context.Context, in *PublishMessagesRequest, opts ...grpc.CallOption) (*PublishMessagesResponse, error)
	PrivateMessages(ctx context.Context, in *PrivateMessagesRequest, opts ...grpc.CallOption) (*PrivateMessagesResponse, error)
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
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetServersResponse)
	err := c.cc.Invoke(ctx, Messages_GetServers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) SaveMessage(ctx context.Context, in *SaveMessageRequest, opts ...grpc.CallOption) (*SaveMessageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SaveMessageResponse)
	err := c.cc.Invoke(ctx, Messages_SaveMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) DeleteMessage(ctx context.Context, in *DeleteMessageRequest, opts ...grpc.CallOption) (*DeleteMessageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteMessageResponse)
	err := c.cc.Invoke(ctx, Messages_DeleteMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) DeleteMessages(ctx context.Context, in *DeleteMessagesRequest, opts ...grpc.CallOption) (*DeleteMessagesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteMessagesResponse)
	err := c.cc.Invoke(ctx, Messages_DeleteMessages_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) DeleteAllUserMessages(ctx context.Context, in *DeleteAllUserMessagesRequest, opts ...grpc.CallOption) (*DeleteAllUserMessagesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteAllUserMessagesResponse)
	err := c.cc.Invoke(ctx, Messages_DeleteAllUserMessages_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) PublishMessages(ctx context.Context, in *PublishMessagesRequest, opts ...grpc.CallOption) (*PublishMessagesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PublishMessagesResponse)
	err := c.cc.Invoke(ctx, Messages_PublishMessages_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) PrivateMessages(ctx context.Context, in *PrivateMessagesRequest, opts ...grpc.CallOption) (*PrivateMessagesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PrivateMessagesResponse)
	err := c.cc.Invoke(ctx, Messages_PrivateMessages_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) UpdateMessage(ctx context.Context, in *UpdateMessageRequest, opts ...grpc.CallOption) (*UpdateMessageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMessageResponse)
	err := c.cc.Invoke(ctx, Messages_UpdateMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) ReadOneMessage(ctx context.Context, in *ReadOneMessageRequest, opts ...grpc.CallOption) (*Message, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Message)
	err := c.cc.Invoke(ctx, Messages_ReadOneMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) ReadAllMessages(ctx context.Context, in *ReadMessagesRequest, opts ...grpc.CallOption) (*ReadMessagesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadMessagesResponse)
	err := c.cc.Invoke(ctx, Messages_ReadAllMessages_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messagesClient) ReadThreadMessages(ctx context.Context, in *ReadThreadMessagesRequest, opts ...grpc.CallOption) (*ReadThreadMessagesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadThreadMessagesResponse)
	err := c.cc.Invoke(ctx, Messages_ReadThreadMessages_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessagesServer is the server API for Messages service.
// All implementations must embed UnimplementedMessagesServer
// for forward compatibility.
type MessagesServer interface {
	GetServers(context.Context, *GetServersRequest) (*GetServersResponse, error)
	SaveMessage(context.Context, *SaveMessageRequest) (*SaveMessageResponse, error)
	DeleteMessage(context.Context, *DeleteMessageRequest) (*DeleteMessageResponse, error)
	DeleteMessages(context.Context, *DeleteMessagesRequest) (*DeleteMessagesResponse, error)
	DeleteAllUserMessages(context.Context, *DeleteAllUserMessagesRequest) (*DeleteAllUserMessagesResponse, error)
	PublishMessages(context.Context, *PublishMessagesRequest) (*PublishMessagesResponse, error)
	PrivateMessages(context.Context, *PrivateMessagesRequest) (*PrivateMessagesResponse, error)
	UpdateMessage(context.Context, *UpdateMessageRequest) (*UpdateMessageResponse, error)
	ReadOneMessage(context.Context, *ReadOneMessageRequest) (*Message, error)
	ReadAllMessages(context.Context, *ReadMessagesRequest) (*ReadMessagesResponse, error)
	ReadThreadMessages(context.Context, *ReadThreadMessagesRequest) (*ReadThreadMessagesResponse, error)
	mustEmbedUnimplementedMessagesServer()
}

// UnimplementedMessagesServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMessagesServer struct{}

func (UnimplementedMessagesServer) GetServers(context.Context, *GetServersRequest) (*GetServersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetServers not implemented")
}
func (UnimplementedMessagesServer) SaveMessage(context.Context, *SaveMessageRequest) (*SaveMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveMessage not implemented")
}
func (UnimplementedMessagesServer) DeleteMessage(context.Context, *DeleteMessageRequest) (*DeleteMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMessage not implemented")
}
func (UnimplementedMessagesServer) DeleteMessages(context.Context, *DeleteMessagesRequest) (*DeleteMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMessages not implemented")
}
func (UnimplementedMessagesServer) DeleteAllUserMessages(context.Context, *DeleteAllUserMessagesRequest) (*DeleteAllUserMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAllUserMessages not implemented")
}
func (UnimplementedMessagesServer) PublishMessages(context.Context, *PublishMessagesRequest) (*PublishMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishMessages not implemented")
}
func (UnimplementedMessagesServer) PrivateMessages(context.Context, *PrivateMessagesRequest) (*PrivateMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrivateMessages not implemented")
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
func (UnimplementedMessagesServer) testEmbeddedByValue()                  {}

// UnsafeMessagesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessagesServer will
// result in compilation errors.
type UnsafeMessagesServer interface {
	mustEmbedUnimplementedMessagesServer()
}

func RegisterMessagesServer(s grpc.ServiceRegistrar, srv MessagesServer) {
	// If the following call pancis, it indicates UnimplementedMessagesServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Messages_ServiceDesc, srv)
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
		FullMethod: Messages_GetServers_FullMethodName,
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
		FullMethod: Messages_SaveMessage_FullMethodName,
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
		FullMethod: Messages_DeleteMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagesServer).DeleteMessage(ctx, req.(*DeleteMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messages_DeleteMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagesServer).DeleteMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Messages_DeleteMessages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagesServer).DeleteMessages(ctx, req.(*DeleteMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messages_DeleteAllUserMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAllUserMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagesServer).DeleteAllUserMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Messages_DeleteAllUserMessages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagesServer).DeleteAllUserMessages(ctx, req.(*DeleteAllUserMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messages_PublishMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagesServer).PublishMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Messages_PublishMessages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagesServer).PublishMessages(ctx, req.(*PublishMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messages_PrivateMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrivateMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagesServer).PrivateMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Messages_PrivateMessages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagesServer).PrivateMessages(ctx, req.(*PrivateMessagesRequest))
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
		FullMethod: Messages_UpdateMessage_FullMethodName,
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
		FullMethod: Messages_ReadOneMessage_FullMethodName,
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
		FullMethod: Messages_ReadAllMessages_FullMethodName,
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
		FullMethod: Messages_ReadThreadMessages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagesServer).ReadThreadMessages(ctx, req.(*ReadThreadMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Messages_ServiceDesc is the grpc.ServiceDesc for Messages service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Messages_ServiceDesc = grpc.ServiceDesc{
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
			MethodName: "DeleteMessages",
			Handler:    _Messages_DeleteMessages_Handler,
		},
		{
			MethodName: "DeleteAllUserMessages",
			Handler:    _Messages_DeleteAllUserMessages_Handler,
		},
		{
			MethodName: "PublishMessages",
			Handler:    _Messages_PublishMessages_Handler,
		},
		{
			MethodName: "PrivateMessages",
			Handler:    _Messages_PrivateMessages_Handler,
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
