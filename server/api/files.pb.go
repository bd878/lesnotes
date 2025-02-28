// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: protos/files.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	CreateUtcNano int64  `protobuf:"varint,2,opt,name=create_utc_nano,json=createUtcNano,proto3" json:"create_utc_nano,omitempty"`
	Id            int32  `protobuf:"varint,3,opt,name=id,proto3" json:"id,omitempty"`
	UserId        int32  `protobuf:"varint,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Error         string `protobuf:"bytes,5,opt,name=error,proto3" json:"error,omitempty"`
	Size          int64  `protobuf:"varint,6,opt,name=size,proto3" json:"size,omitempty"`
}

func (x *File) Reset() {
	*x = File{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_files_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_protos_files_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_protos_files_proto_rawDescGZIP(), []int{0}
}

func (x *File) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *File) GetCreateUtcNano() int64 {
	if x != nil {
		return x.CreateUtcNano
	}
	return 0
}

func (x *File) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *File) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *File) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *File) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

type FileChunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chunk []byte `protobuf:"bytes,1,opt,name=chunk,proto3" json:"chunk,omitempty"`
}

func (x *FileChunk) Reset() {
	*x = FileChunk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_files_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileChunk) ProtoMessage() {}

func (x *FileChunk) ProtoReflect() protoreflect.Message {
	mi := &file_protos_files_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileChunk.ProtoReflect.Descriptor instead.
func (*FileChunk) Descriptor() ([]byte, []int) {
	return file_protos_files_proto_rawDescGZIP(), []int{1}
}

func (x *FileChunk) GetChunk() []byte {
	if x != nil {
		return x.Chunk
	}
	return nil
}

type FileData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//	*FileData_File
	//	*FileData_Chunk
	Data isFileData_Data `protobuf_oneof:"data"`
}

func (x *FileData) Reset() {
	*x = FileData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_files_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileData) ProtoMessage() {}

func (x *FileData) ProtoReflect() protoreflect.Message {
	mi := &file_protos_files_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileData.ProtoReflect.Descriptor instead.
func (*FileData) Descriptor() ([]byte, []int) {
	return file_protos_files_proto_rawDescGZIP(), []int{2}
}

func (m *FileData) GetData() isFileData_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *FileData) GetFile() *File {
	if x, ok := x.GetData().(*FileData_File); ok {
		return x.File
	}
	return nil
}

func (x *FileData) GetChunk() []byte {
	if x, ok := x.GetData().(*FileData_Chunk); ok {
		return x.Chunk
	}
	return nil
}

type isFileData_Data interface {
	isFileData_Data()
}

type FileData_File struct {
	File *File `protobuf:"bytes,1,opt,name=file,proto3,oneof"`
}

type FileData_Chunk struct {
	Chunk []byte `protobuf:"bytes,2,opt,name=chunk,proto3,oneof"`
}

func (*FileData_File) isFileData_Data() {}

func (*FileData_Chunk) isFileData_Data() {}

type SaveFileStreamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	File *File `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
}

func (x *SaveFileStreamResponse) Reset() {
	*x = SaveFileStreamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_files_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveFileStreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveFileStreamResponse) ProtoMessage() {}

func (x *SaveFileStreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_files_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveFileStreamResponse.ProtoReflect.Descriptor instead.
func (*SaveFileStreamResponse) Descriptor() ([]byte, []int) {
	return file_protos_files_proto_rawDescGZIP(), []int{3}
}

func (x *SaveFileStreamResponse) GetFile() *File {
	if x != nil {
		return x.File
	}
	return nil
}

type ReadFileStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId int32 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *ReadFileStreamRequest) Reset() {
	*x = ReadFileStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_files_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadFileStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadFileStreamRequest) ProtoMessage() {}

func (x *ReadFileStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_files_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadFileStreamRequest.ProtoReflect.Descriptor instead.
func (*ReadFileStreamRequest) Descriptor() ([]byte, []int) {
	return file_protos_files_proto_rawDescGZIP(), []int{4}
}

func (x *ReadFileStreamRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ReadFileStreamRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type ReadBatchFilesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int32   `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Ids    []int32 `protobuf:"varint,2,rep,packed,name=ids,proto3" json:"ids,omitempty"`
}

func (x *ReadBatchFilesRequest) Reset() {
	*x = ReadBatchFilesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_files_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadBatchFilesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadBatchFilesRequest) ProtoMessage() {}

func (x *ReadBatchFilesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_files_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadBatchFilesRequest.ProtoReflect.Descriptor instead.
func (*ReadBatchFilesRequest) Descriptor() ([]byte, []int) {
	return file_protos_files_proto_rawDescGZIP(), []int{5}
}

func (x *ReadBatchFilesRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ReadBatchFilesRequest) GetIds() []int32 {
	if x != nil {
		return x.Ids
	}
	return nil
}

type ReadBatchFilesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Files map[int32]*File `protobuf:"bytes,1,rep,name=files,proto3" json:"files,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ReadBatchFilesResponse) Reset() {
	*x = ReadBatchFilesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_files_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadBatchFilesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadBatchFilesResponse) ProtoMessage() {}

func (x *ReadBatchFilesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_files_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadBatchFilesResponse.ProtoReflect.Descriptor instead.
func (*ReadBatchFilesResponse) Descriptor() ([]byte, []int) {
	return file_protos_files_proto_rawDescGZIP(), []int{6}
}

func (x *ReadBatchFilesResponse) GetFiles() map[int32]*File {
	if x != nil {
		return x.Files
	}
	return nil
}

var File_protos_files_proto protoreflect.FileDescriptor

var file_protos_files_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x22, 0x95,
	0x01, 0x0a, 0x04, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0f, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x75, 0x74, 0x63, 0x5f, 0x6e, 0x61, 0x6e, 0x6f, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x74, 0x63, 0x4e,
	0x61, 0x6e, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x21, 0x0a, 0x09, 0x46, 0x69, 0x6c, 0x65, 0x43, 0x68,
	0x75, 0x6e, 0x6b, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x22, 0x50, 0x0a, 0x08, 0x46, 0x69, 0x6c,
	0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x24, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x48, 0x00, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x05, 0x63,
	0x68, 0x75, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x05, 0x63, 0x68,
	0x75, 0x6e, 0x6b, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x3c, 0x0a, 0x16, 0x53,
	0x61, 0x76, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x40, 0x0a, 0x15, 0x52, 0x65, 0x61,
	0x64, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x42, 0x0a, 0x15, 0x52,
	0x65, 0x61, 0x64, 0x42, 0x61, 0x74, 0x63, 0x68, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x10, 0x0a,
	0x03, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x05, 0x52, 0x03, 0x69, 0x64, 0x73, 0x22,
	0xa5, 0x01, 0x0a, 0x16, 0x52, 0x65, 0x61, 0x64, 0x42, 0x61, 0x74, 0x63, 0x68, 0x46, 0x69, 0x6c,
	0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x05, 0x66, 0x69,
	0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x66, 0x69, 0x6c, 0x65,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x42, 0x61, 0x74, 0x63, 0x68, 0x46, 0x69,
	0x6c, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x1a, 0x48, 0x0a,
	0x0a, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x24, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x66,
	0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0xf5, 0x01, 0x0a, 0x05, 0x46, 0x69, 0x6c, 0x65,
	0x73, 0x12, 0x55, 0x0a, 0x0e, 0x52, 0x65, 0x61, 0x64, 0x42, 0x61, 0x74, 0x63, 0x68, 0x46, 0x69,
	0x6c, 0x65, 0x73, 0x12, 0x1f, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52,
	0x65, 0x61, 0x64, 0x42, 0x61, 0x74, 0x63, 0x68, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x52, 0x65, 0x61, 0x64, 0x42, 0x61, 0x74, 0x63, 0x68, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4a, 0x0a, 0x0e, 0x53, 0x61, 0x76, 0x65,
	0x46, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x12, 0x2e, 0x66, 0x69, 0x6c,
	0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x20,
	0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x46, 0x69,
	0x6c, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x28, 0x01, 0x12, 0x49, 0x0a, 0x0e, 0x52, 0x65, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x1f, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x22, 0x00, 0x30, 0x01, 0x42,
	0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x64,
	0x38, 0x37, 0x38, 0x2f, 0x67, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x79, 0x2f, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_files_proto_rawDescOnce sync.Once
	file_protos_files_proto_rawDescData = file_protos_files_proto_rawDesc
)

func file_protos_files_proto_rawDescGZIP() []byte {
	file_protos_files_proto_rawDescOnce.Do(func() {
		file_protos_files_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_files_proto_rawDescData)
	})
	return file_protos_files_proto_rawDescData
}

var file_protos_files_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_protos_files_proto_goTypes = []interface{}{
	(*File)(nil),                   // 0: files.v1.File
	(*FileChunk)(nil),              // 1: files.v1.FileChunk
	(*FileData)(nil),               // 2: files.v1.FileData
	(*SaveFileStreamResponse)(nil), // 3: files.v1.SaveFileStreamResponse
	(*ReadFileStreamRequest)(nil),  // 4: files.v1.ReadFileStreamRequest
	(*ReadBatchFilesRequest)(nil),  // 5: files.v1.ReadBatchFilesRequest
	(*ReadBatchFilesResponse)(nil), // 6: files.v1.ReadBatchFilesResponse
	nil,                            // 7: files.v1.ReadBatchFilesResponse.FilesEntry
}
var file_protos_files_proto_depIdxs = []int32{
	0, // 0: files.v1.FileData.file:type_name -> files.v1.File
	0, // 1: files.v1.SaveFileStreamResponse.file:type_name -> files.v1.File
	7, // 2: files.v1.ReadBatchFilesResponse.files:type_name -> files.v1.ReadBatchFilesResponse.FilesEntry
	0, // 3: files.v1.ReadBatchFilesResponse.FilesEntry.value:type_name -> files.v1.File
	5, // 4: files.v1.Files.ReadBatchFiles:input_type -> files.v1.ReadBatchFilesRequest
	2, // 5: files.v1.Files.SaveFileStream:input_type -> files.v1.FileData
	4, // 6: files.v1.Files.ReadFileStream:input_type -> files.v1.ReadFileStreamRequest
	6, // 7: files.v1.Files.ReadBatchFiles:output_type -> files.v1.ReadBatchFilesResponse
	3, // 8: files.v1.Files.SaveFileStream:output_type -> files.v1.SaveFileStreamResponse
	2, // 9: files.v1.Files.ReadFileStream:output_type -> files.v1.FileData
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_protos_files_proto_init() }
func file_protos_files_proto_init() {
	if File_protos_files_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_files_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*File); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_files_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileChunk); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_files_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_files_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveFileStreamResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_files_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadFileStreamRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_files_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadBatchFilesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_files_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadBatchFilesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_protos_files_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*FileData_File)(nil),
		(*FileData_Chunk)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_files_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_files_proto_goTypes,
		DependencyIndexes: file_protos_files_proto_depIdxs,
		MessageInfos:      file_protos_files_proto_msgTypes,
	}.Build()
	File_protos_files_proto = out.File
	file_protos_files_proto_rawDesc = nil
	file_protos_files_proto_goTypes = nil
	file_protos_files_proto_depIdxs = nil
}
