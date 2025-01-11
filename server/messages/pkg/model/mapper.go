package model

import "github.com/bd878/gallery/server/api"

func MessageFromProto(proto *api.Message) *Message {
  return &Message{
    ID:                proto.Id,
    CreateUTCNano:     proto.CreateUtcNano,
    UpdateUTCNano:     proto.UpdateUtcNano,
    UserID:            proto.UserId,
    Text:              proto.Text,
    FileID:            proto.FileId,
  }
}

func MessageToProto(msg *Message) *api.Message {
  return &api.Message{
    Id:                msg.ID,
    UserId:            msg.UserID,
    CreateUtcNano:     msg.CreateUTCNano,
    UpdateUtcNano:     msg.UpdateUTCNano,
    Text:              msg.Text,
    FileId:            msg.FileID,
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