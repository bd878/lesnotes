package model

import "github.com/bd878/gallery/server/api"

type MessageId int

type FileId string

// This message handler passes to repository
type Message struct {
  Id MessageId `json:"id"`
  CreateTime string `json:"createtime"`
  UserId int `json:"userid"`
  Value string `json:"value"`
  FileName string `json:"filename"`
  FileId FileId `json:"fileid"`
  LogIndex uint64 `json:"logindex,omitempty"`
  LogTerm uint64 `json:"logterm,omitempty"`
}

// Response to return to the client
type ServerResponse struct {
  Status string `json:"status"`
  Description string `json:"description"`
}

type NewMessageServerResponse struct {
  ServerResponse
  Message Message `json:"message"`
}

func ProtoToMessage(proto api.Message) Message {
  return Message{
    Id:          MessageId(proto.Id),
    CreateTime:  proto.CreateTime,
    UserId:      int(proto.UserId),
    Value:       string(proto.Value),
    FileName:    proto.FileName,
    FileId:      FileId(proto.FileId),
  }
}

func MessageToProto(msg Message) api.Message {
  return api.Message{
    Id:          uint32(msg.Id),
    CreateTime:  msg.CreateTime,
    Value:       []byte(msg.Value),
    FileName:    msg.FileName,
    FileId:      string(msg.FileId),
  }
}