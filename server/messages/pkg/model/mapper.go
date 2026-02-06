package model

import (
	"github.com/bd878/gallery/server/api"
)

func MessageFromProto(proto *api.Message) *Message {
	return &Message{
		ID:             proto.Id,
		CreateUTCNano:  proto.CreateUtcNano,
		UpdateUTCNano:  proto.UpdateUtcNano,
		UserID:         proto.UserId,
		Text:           proto.Text,
		Title:          proto.Title,
		FileIDs:        proto.FileIds,
		Private:        proto.Private,
		Name:           proto.Name,
		Translations:   MapTranslationsFromProto(TranslationFromProto, proto.Translations),
	}
}

func MessageToProto(msg *Message) *api.Message {
	return &api.Message{
		Id:             msg.ID,
		UserId:         msg.UserID,
		CreateUtcNano:  msg.CreateUTCNano,
		UpdateUtcNano:  msg.UpdateUTCNano,
		Text:           msg.Text,
		Title:          msg.Title,
		FileIds:        msg.FileIDs,
		Private:        msg.Private,
		Name:           msg.Name,
		Translations:   MapTranslationsToProto(TranslationToProto, msg.Translations),
	}
}

func TranslationFromProto(proto *api.Translation) *Translation {
	return &Translation{
		MessageID:      proto.MessageId,
		Lang:           proto.Lang,
		Title:          proto.Title,
		Text:           proto.Text,
	}
}

func TranslationToProto(translation *Translation) *api.Translation {
	return &api.Translation{
		MessageId:      translation.MessageID,
		Lang:           translation.Lang,
		Title:          translation.Title,
		Text:           translation.Text,
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

func MapTranslationsToProto(mapper (func(*Translation) *api.Translation), translations []*Translation) []*api.Translation {
	res := make([]*api.Translation, len(translations))
	for i, translation := range translations {
		res[i] = mapper(translation)
	}
	return res
}

func MapTranslationsFromProto(mapper (func(*api.Translation) *Translation), translations []*api.Translation) []*Translation {
	res := make([]*Translation, len(translations))
	for i, translation := range translations {
		res[i] = mapper(translation)
	}
	return res
}
