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
		CreatedAt:      proto.CreatedAt,
		UpdatedAt:      proto.UpdatedAt,
	}
}

func TranslationToProto(translation *Translation) *api.Translation {
	return &api.Translation{
		MessageId:      translation.MessageID,
		Lang:           translation.Lang,
		Title:          translation.Title,
		Text:           translation.Text,
		CreatedAt:      translation.CreatedAt,
		UpdatedAt:      translation.UpdatedAt,
	}
}

func TranslationPreviewFromProto(proto *api.TranslationPreview) *TranslationPreview {
	return &TranslationPreview{
		MessageID:      proto.MessageId,
		Lang:           proto.Lang,
		Title:          proto.Title,
		CreatedAt:      proto.CreatedAt,
		UpdatedAt:      proto.UpdatedAt,
	}
}

func TranslationPreviewToProto(preview *TranslationPreview) *api.TranslationPreview {
	return &api.TranslationPreview{
		MessageId:      preview.MessageID,
		Lang:           preview.Lang,
		Title:          preview.Title,
		CreatedAt:      preview.CreatedAt,
		UpdatedAt:      preview.UpdatedAt,
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

func MapTranslationPreviewsToProto(mapper (func(*TranslationPreview) *api.TranslationPreview), previews []*TranslationPreview) []*api.TranslationPreview {
	res := make([]*api.TranslationPreview, len(previews))
	for i, preview := range previews {
		res[i] = mapper(preview)
	}
	return res
}

func MapTranslationPreviewsFromProto(mapper (func(*api.TranslationPreview) *TranslationPreview), previews []*api.TranslationPreview) []*TranslationPreview {
	res := make([]*TranslationPreview, len(previews))
	for i, preview := range previews {
		res[i] = mapper(preview)
	}
	return res
}