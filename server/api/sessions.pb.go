// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.19.4
// source: protos/sessions.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Session struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	UserId         int32                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Token          string                 `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	ExpiresUtcNano int64                  `protobuf:"varint,4,opt,name=expires_utc_nano,json=expiresUtcNano,proto3" json:"expires_utc_nano,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *Session) Reset() {
	*x = Session{}
	mi := &file_protos_sessions_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Session) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Session) ProtoMessage() {}

func (x *Session) ProtoReflect() protoreflect.Message {
	mi := &file_protos_sessions_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Session.ProtoReflect.Descriptor instead.
func (*Session) Descriptor() ([]byte, []int) {
	return file_protos_sessions_proto_rawDescGZIP(), []int{0}
}

func (x *Session) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Session) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *Session) GetExpiresUtcNano() int64 {
	if x != nil {
		return x.ExpiresUtcNano
	}
	return 0
}

type CreateSessionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int32                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateSessionRequest) Reset() {
	*x = CreateSessionRequest{}
	mi := &file_protos_sessions_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSessionRequest) ProtoMessage() {}

func (x *CreateSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_sessions_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSessionRequest.ProtoReflect.Descriptor instead.
func (*CreateSessionRequest) Descriptor() ([]byte, []int) {
	return file_protos_sessions_proto_rawDescGZIP(), []int{1}
}

func (x *CreateSessionRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type ListUserSessionsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int32                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListUserSessionsRequest) Reset() {
	*x = ListUserSessionsRequest{}
	mi := &file_protos_sessions_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListUserSessionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUserSessionsRequest) ProtoMessage() {}

func (x *ListUserSessionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_sessions_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListUserSessionsRequest.ProtoReflect.Descriptor instead.
func (*ListUserSessionsRequest) Descriptor() ([]byte, []int) {
	return file_protos_sessions_proto_rawDescGZIP(), []int{2}
}

func (x *ListUserSessionsRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type ListUserSessionsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Sessions      []*Session             `protobuf:"bytes,1,rep,name=sessions,proto3" json:"sessions,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListUserSessionsResponse) Reset() {
	*x = ListUserSessionsResponse{}
	mi := &file_protos_sessions_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListUserSessionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUserSessionsResponse) ProtoMessage() {}

func (x *ListUserSessionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_sessions_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListUserSessionsResponse.ProtoReflect.Descriptor instead.
func (*ListUserSessionsResponse) Descriptor() ([]byte, []int) {
	return file_protos_sessions_proto_rawDescGZIP(), []int{3}
}

func (x *ListUserSessionsResponse) GetSessions() []*Session {
	if x != nil {
		return x.Sessions
	}
	return nil
}

type GetSessionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetSessionRequest) Reset() {
	*x = GetSessionRequest{}
	mi := &file_protos_sessions_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSessionRequest) ProtoMessage() {}

func (x *GetSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_sessions_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSessionRequest.ProtoReflect.Descriptor instead.
func (*GetSessionRequest) Descriptor() ([]byte, []int) {
	return file_protos_sessions_proto_rawDescGZIP(), []int{4}
}

func (x *GetSessionRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type RemoveSessionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RemoveSessionRequest) Reset() {
	*x = RemoveSessionRequest{}
	mi := &file_protos_sessions_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RemoveSessionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveSessionRequest) ProtoMessage() {}

func (x *RemoveSessionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_sessions_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveSessionRequest.ProtoReflect.Descriptor instead.
func (*RemoveSessionRequest) Descriptor() ([]byte, []int) {
	return file_protos_sessions_proto_rawDescGZIP(), []int{5}
}

func (x *RemoveSessionRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type RemoveAllSessionsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int32                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RemoveAllSessionsRequest) Reset() {
	*x = RemoveAllSessionsRequest{}
	mi := &file_protos_sessions_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RemoveAllSessionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveAllSessionsRequest) ProtoMessage() {}

func (x *RemoveAllSessionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_sessions_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveAllSessionsRequest.ProtoReflect.Descriptor instead.
func (*RemoveAllSessionsRequest) Descriptor() ([]byte, []int) {
	return file_protos_sessions_proto_rawDescGZIP(), []int{6}
}

func (x *RemoveAllSessionsRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type RemoveAllSessionsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RemoveAllSessionsResponse) Reset() {
	*x = RemoveAllSessionsResponse{}
	mi := &file_protos_sessions_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RemoveAllSessionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveAllSessionsResponse) ProtoMessage() {}

func (x *RemoveAllSessionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_sessions_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveAllSessionsResponse.ProtoReflect.Descriptor instead.
func (*RemoveAllSessionsResponse) Descriptor() ([]byte, []int) {
	return file_protos_sessions_proto_rawDescGZIP(), []int{7}
}

type AddSessionResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddSessionResponse) Reset() {
	*x = AddSessionResponse{}
	mi := &file_protos_sessions_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddSessionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddSessionResponse) ProtoMessage() {}

func (x *AddSessionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_sessions_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddSessionResponse.ProtoReflect.Descriptor instead.
func (*AddSessionResponse) Descriptor() ([]byte, []int) {
	return file_protos_sessions_proto_rawDescGZIP(), []int{8}
}

type RemoveSessionResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RemoveSessionResponse) Reset() {
	*x = RemoveSessionResponse{}
	mi := &file_protos_sessions_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RemoveSessionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveSessionResponse) ProtoMessage() {}

func (x *RemoveSessionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_sessions_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveSessionResponse.ProtoReflect.Descriptor instead.
func (*RemoveSessionResponse) Descriptor() ([]byte, []int) {
	return file_protos_sessions_proto_rawDescGZIP(), []int{9}
}

var File_protos_sessions_proto protoreflect.FileDescriptor

var file_protos_sessions_proto_rawDesc = string([]byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x76, 0x31, 0x22, 0x62, 0x0a, 0x07, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x28,
	0x0a, 0x10, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x75, 0x74, 0x63, 0x5f, 0x6e, 0x61,
	0x6e, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65,
	0x73, 0x55, 0x74, 0x63, 0x4e, 0x61, 0x6e, 0x6f, 0x22, 0x2f, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x32, 0x0a, 0x17, 0x4c, 0x69, 0x73,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x4c, 0x0a,
	0x18, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x08, 0x73, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x08, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x29, 0x0a, 0x11, 0x47,
	0x65, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x2c, 0x0a, 0x14, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x33, 0x0a, 0x18, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x41, 0x6c,
	0x6c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x1b, 0x0a, 0x19, 0x52, 0x65, 0x6d,
	0x6f, 0x76, 0x65, 0x41, 0x6c, 0x6c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x14, 0x0a, 0x12, 0x41, 0x64, 0x64, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x17, 0x0a, 0x15,
	0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x8c, 0x03, 0x0a, 0x08, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x53, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x24, 0x2e, 0x73, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x25, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1e,
	0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14,
	0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x41, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x21,
	0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x14, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x4f, 0x0a, 0x06, 0x52, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x12, 0x21, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5a, 0x0a, 0x09, 0x52, 0x65, 0x6d, 0x6f,
	0x76, 0x65, 0x41, 0x6c, 0x6c, 0x12, 0x25, 0x2e, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x41, 0x6c, 0x6c, 0x53, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x73,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x41, 0x6c, 0x6c, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x62, 0x64, 0x38, 0x37, 0x38, 0x2f, 0x67, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x79,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
})

var (
	file_protos_sessions_proto_rawDescOnce sync.Once
	file_protos_sessions_proto_rawDescData []byte
)

func file_protos_sessions_proto_rawDescGZIP() []byte {
	file_protos_sessions_proto_rawDescOnce.Do(func() {
		file_protos_sessions_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_protos_sessions_proto_rawDesc), len(file_protos_sessions_proto_rawDesc)))
	})
	return file_protos_sessions_proto_rawDescData
}

var file_protos_sessions_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_protos_sessions_proto_goTypes = []any{
	(*Session)(nil),                   // 0: sessions.v1.Session
	(*CreateSessionRequest)(nil),      // 1: sessions.v1.CreateSessionRequest
	(*ListUserSessionsRequest)(nil),   // 2: sessions.v1.ListUserSessionsRequest
	(*ListUserSessionsResponse)(nil),  // 3: sessions.v1.ListUserSessionsResponse
	(*GetSessionRequest)(nil),         // 4: sessions.v1.GetSessionRequest
	(*RemoveSessionRequest)(nil),      // 5: sessions.v1.RemoveSessionRequest
	(*RemoveAllSessionsRequest)(nil),  // 6: sessions.v1.RemoveAllSessionsRequest
	(*RemoveAllSessionsResponse)(nil), // 7: sessions.v1.RemoveAllSessionsResponse
	(*AddSessionResponse)(nil),        // 8: sessions.v1.AddSessionResponse
	(*RemoveSessionResponse)(nil),     // 9: sessions.v1.RemoveSessionResponse
}
var file_protos_sessions_proto_depIdxs = []int32{
	0, // 0: sessions.v1.ListUserSessionsResponse.sessions:type_name -> sessions.v1.Session
	2, // 1: sessions.v1.Sessions.List:input_type -> sessions.v1.ListUserSessionsRequest
	4, // 2: sessions.v1.Sessions.Get:input_type -> sessions.v1.GetSessionRequest
	1, // 3: sessions.v1.Sessions.Create:input_type -> sessions.v1.CreateSessionRequest
	5, // 4: sessions.v1.Sessions.Remove:input_type -> sessions.v1.RemoveSessionRequest
	6, // 5: sessions.v1.Sessions.RemoveAll:input_type -> sessions.v1.RemoveAllSessionsRequest
	3, // 6: sessions.v1.Sessions.List:output_type -> sessions.v1.ListUserSessionsResponse
	0, // 7: sessions.v1.Sessions.Get:output_type -> sessions.v1.Session
	0, // 8: sessions.v1.Sessions.Create:output_type -> sessions.v1.Session
	9, // 9: sessions.v1.Sessions.Remove:output_type -> sessions.v1.RemoveSessionResponse
	7, // 10: sessions.v1.Sessions.RemoveAll:output_type -> sessions.v1.RemoveAllSessionsResponse
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_protos_sessions_proto_init() }
func file_protos_sessions_proto_init() {
	if File_protos_sessions_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_protos_sessions_proto_rawDesc), len(file_protos_sessions_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_sessions_proto_goTypes,
		DependencyIndexes: file_protos_sessions_proto_depIdxs,
		MessageInfos:      file_protos_sessions_proto_msgTypes,
	}.Build()
	File_protos_sessions_proto = out.File
	file_protos_sessions_proto_goTypes = nil
	file_protos_sessions_proto_depIdxs = nil
}
