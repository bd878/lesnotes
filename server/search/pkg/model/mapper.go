package model

import "github.com/bd878/gallery/server/api/search"

func MessageFromProto(proto *search.SearchMessage) *Message {
	return &Message{
		ID:             proto.Id,
		UserID:         proto.UserId,
		Text:           proto.Text,
		Title:          proto.Title,
		Private:        proto.Private,
		Name:           proto.Name,
		CreatedAt:      proto.CreatedAt,
		UpdatedAt:      proto.UpdatedAt,
	}
}

func MessageToProto(msg *Message) *search.SearchMessage {
	return &search.SearchMessage{
		Id:             msg.ID,
		UserId:         msg.UserID,
		Text:           msg.Text,
		Title:          msg.Title,
		Private:        msg.Private,
		Name:           msg.Name,
		CreatedAt:      msg.CreatedAt,
		UpdatedAt:      msg.UpdatedAt,
	}
}

func MapMessagesToProto(mapper (func(*Message) *search.SearchMessage), msgs []*Message) []*search.SearchMessage {
	res := make([]*search.SearchMessage, len(msgs))
	for i, msg := range msgs {
		res[i] = mapper(msg)
	}
	return res
}

func MapMessagesFromProto(mapper (func(*search.SearchMessage) *Message), msgs []*search.SearchMessage) []*Message {
	res := make([]*Message, len(msgs))
	for i, msg := range msgs {
		res[i] = mapper(msg)
	}
	return res
}