package model

import (
	"github.com/bd878/gallery/server/api"
)

func MessageFromProto(proto *api.SearchMessage) *Message {
	return &Message{
		ID:             proto.Id,
		UserID:         proto.UserId,
		Text:           proto.Text,
		Title:          proto.Title,
		Private:        proto.Private,
		Name:           proto.Name,
	}
}

func MessageToProto(msg *Message) *api.SearchMessage {
	return &api.SearchMessage{
		Id:             msg.ID,
		UserId:         msg.UserID,
		Text:           msg.Text,
		Title:          msg.Title,
		Private:        msg.Private,
		Name:           msg.Name,
	}
}

func MapMessagesToProto(mapper (func(*Message) *api.SearchMessage), msgs []*Message) []*api.SearchMessage {
	res := make([]*api.SearchMessage, len(msgs))
	for i, msg := range msgs {
		res[i] = mapper(msg)
	}
	return res
}

func MapMessagesFromProto(mapper (func(*api.SearchMessage) *Message), msgs []*api.SearchMessage) []*Message {
	res := make([]*Message, len(msgs))
	for i, msg := range msgs {
		res[i] = mapper(msg)
	}
	return res
}