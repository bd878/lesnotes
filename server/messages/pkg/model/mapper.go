package model

import "github.com/bd878/gallery/server/api"

func MessageFromProto(proto *api.Message) *Message {
  return &Message{
    Id:             MessageId(proto.Id),
    CreateTime:     proto.CreateTime,
    UpdateUtcNano:  proto.UpdateUtcNano,
    UserId:         int(proto.UserId),
    Value:          string(proto.Value),
    FileName:       proto.FileName,
    FileId:         FileId(proto.FileId),
  }
}

func MessageToProto(msg *Message) *api.Message {
  return &api.Message{
    Id:               uint32(msg.Id),
    UserId:           uint32(msg.UserId),
    CreateTime:       msg.CreateTime,
    UpdateUtcNano:    msg.UpdateUtcNano,
    Value:            []byte(msg.Value),
    FileName:         msg.FileName,
    FileId:           string(msg.FileId),
  }
}

func MapMessagesToProto(mapper (func(*Message) *api.Message), msgs []*Message) []*api.Message {
  res := make([]*api.Message, len(msgs))
  for i, msg := range msgs {
    res[i] = mapper(msg)
  }
  return res
}

func MapMessagesFromProto(mapper (func(*api.Message) *Message), msgs []*api.Message) []*Message {
  res := make([]*Message, len(msgs))
  for i, msg := range msgs {
    res[i] = mapper(msg)
  }
  return res
}