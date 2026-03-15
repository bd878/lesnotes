package model

import (
	"github.com/bd878/gallery/server/api"
)

func MessageFromProto(proto *api.Message) *Message {
	return &Message{
		ID:             proto.Id,
		CreatedAt:      proto.CreatedAt,
		UpdatedAt:      proto.UpdatedAt,
		UserID:         proto.UserId,
		Text:           proto.Text,
		Title:          proto.Title,
		FileIDs:        proto.FileIds,
		Private:        proto.Private,
		Name:           proto.Name,
		Translations:   MapTranslationPreviewsFromProto(TranslationPreviewFromProto, proto.Translations),
	}
}

func MessageToProto(msg *Message) *api.Message {
	return &api.Message{
		Id:             msg.ID,
		UserId:         msg.UserID,
		CreatedAt:      msg.CreatedAt,
		UpdatedAt:      msg.UpdatedAt,
		Text:           msg.Text,
		Title:          msg.Title,
		FileIds:        msg.FileIDs,
		Private:        msg.Private,
		Name:           msg.Name,
		Translations:   MapTranslationPreviewsToProto(TranslationPreviewToProto, msg.Translations),
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

func CommentFromProto(proto *api.Comment) *Comment {
	return &Comment{
		MessageID:      proto.MessageId,
		UserID:         proto.UserId,
		ID:             proto.Id,
		Text:           proto.Text,
		Metadata:       proto.Metadata,
		CreatedAt:      proto.CreatedAt,
		UpdatedAt:      proto.UpdatedAt,
	}
}

func CommentToProto(comment *Comment) *api.Comment {
	return &api.Comment{
		MessageId:     comment.MessageID,
		UserId:        comment.UserID,
		Id:            comment.ID,
		Text:          comment.Text,
		Metadata:      comment.Metadata,
		CreatedAt:     comment.CreatedAt,
		UpdatedAt:     comment.UpdatedAt,
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

func MapCommentsFromProto(mapper (func(*api.Comment) *Comment), comments []*api.Comment) []*Comment {
	res := make([]*Comment, len(comments))
	for i, comment := range comments {
		res[i] = mapper(comment)
	}
	return res
}

func MapCommentsToProto(mapper (func(*Comment) *api.Comment), comments []*Comment) []*api.Comment {
	res := make([]*api.Comment, len(comments))
	for i, comment := range comments {
		res[i] = mapper(comment)
	}
	return res
}