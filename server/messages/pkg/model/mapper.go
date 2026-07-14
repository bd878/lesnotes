package model

import (
	"github.com/bd878/gallery/server/api/messages"
	"github.com/bd878/gallery/server/api/translations"
	"github.com/bd878/gallery/server/api/comments"
	files "github.com/bd878/gallery/server/files/pkg/model"
)

func MessageFromProto(proto *messages.Message) *Message {
	return &Message{
		ID:           proto.Id,
		CreatedAt:    proto.CreatedAt,
		UpdatedAt:    proto.UpdatedAt,
		UserID:       proto.UserId,
		Text:         proto.Text,
		Title:        proto.Title,
		FileIDs:      proto.FileIds,
		Files:        files.MapFilesFromProto(files.FileFromProto, proto.Files),
		Private:      proto.Private,
		Name:         proto.Name,
		Translations: MapTranslationPreviewsFromProto(TranslationPreviewFromProto, proto.Translations),
	}
}

func MessageToProto(msg *Message) *messages.Message {
	return &messages.Message{
		Id:           msg.ID,
		UserId:       msg.UserID,
		CreatedAt:    msg.CreatedAt,
		UpdatedAt:    msg.UpdatedAt,
		Text:         msg.Text,
		Title:        msg.Title,
		FileIds:      msg.FileIDs,
		Files:        files.MapFilesToProto(files.FileToProto, msg.Files),
		Private:      msg.Private,
		Name:         msg.Name,
		Translations: MapTranslationPreviewsToProto(TranslationPreviewToProto, msg.Translations),
	}
}

func TranslationFromProto(proto *translations.Translation) *Translation {
	return &Translation{
		MessageID: proto.Id,
		Lang:      proto.Lang,
		Title:     proto.Title,
		Text:      proto.Text,
		CreatedAt: proto.CreatedAt,
		UpdatedAt: proto.UpdatedAt,
	}
}

func TranslationToProto(translation *Translation) *translations.Translation {
	return &translations.Translation{
		Id:        translation.MessageID,
		Lang:      translation.Lang,
		Title:     translation.Title,
		Text:      translation.Text,
		CreatedAt: translation.CreatedAt,
		UpdatedAt: translation.UpdatedAt,
	}
}

func CommentFromProto(proto *comments.Comment) *Comment {
	return &Comment{
		MessageID: proto.MessageId,
		UserID:    proto.UserId,
		ID:        proto.Id,
		Text:      proto.Text,
		Metadata:  proto.Metadata,
		CreatedAt: proto.CreatedAt,
		UpdatedAt: proto.UpdatedAt,
	}
}

func CommentToProto(comment *Comment) *comments.Comment {
	return &comments.Comment{
		MessageId: comment.MessageID,
		UserId:    comment.UserID,
		Id:        comment.ID,
		Text:      comment.Text,
		Metadata:  comment.Metadata,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}

func TranslationPreviewFromProto(proto *translations.TranslationPreview) *TranslationPreview {
	return &TranslationPreview{
		MessageID: proto.Id,
		Lang:      proto.Lang,
		Title:     proto.Title,
		CreatedAt: proto.CreatedAt,
		UpdatedAt: proto.UpdatedAt,
	}
}

func TranslationPreviewToProto(preview *TranslationPreview) *translations.TranslationPreview {
	return &translations.TranslationPreview{
		Id:        preview.MessageID,
		Lang:      preview.Lang,
		Title:     preview.Title,
		CreatedAt: preview.CreatedAt,
		UpdatedAt: preview.UpdatedAt,
	}
}

func MapMessagesToProto(mapper func(*Message) *messages.Message, msgs []*Message) []*messages.Message {
	res := make([]*messages.Message, len(msgs))
	for i, msg := range msgs {
		res[i] = mapper(msg)
	}
	return res
}

func MapMessagesFromProto(mapper func(*messages.Message) *Message, msgs []*messages.Message) []*Message {
	res := make([]*Message, len(msgs))
	for i, msg := range msgs {
		res[i] = mapper(msg)
	}
	return res
}

func MapTranslationsToProto(mapper func(*Translation) *translations.Translation, list []*Translation) []*translations.Translation {
	res := make([]*translations.Translation, len(list))
	for i, translation := range list {
		res[i] = mapper(translation)
	}
	return res
}

func MapTranslationsFromProto(mapper func(*translations.Translation) *Translation, list []*translations.Translation) []*Translation {
	res := make([]*Translation, len(list))
	for i, translation := range list {
		res[i] = mapper(translation)
	}
	return res
}

func MapTranslationPreviewsToProto(mapper func(*TranslationPreview) *translations.TranslationPreview, previews []*TranslationPreview) []*translations.TranslationPreview {
	res := make([]*translations.TranslationPreview, len(previews))
	for i, preview := range previews {
		res[i] = mapper(preview)
	}
	return res
}

func MapTranslationPreviewsFromProto(mapper func(*translations.TranslationPreview) *TranslationPreview, previews []*translations.TranslationPreview) []*TranslationPreview {
	res := make([]*TranslationPreview, len(previews))
	for i, preview := range previews {
		res[i] = mapper(preview)
	}
	return res
}

func MapCommentsFromProto(mapper func(*comments.Comment) *Comment, list []*comments.Comment) []*Comment {
	res := make([]*Comment, len(list))
	for i, comment := range list {
		res[i] = mapper(comment)
	}
	return res
}

func MapCommentsToProto(mapper func(*Comment) *comments.Comment, list []*Comment) []*comments.Comment {
	res := make([]*comments.Comment, len(list))
	for i, comment := range list {
		res[i] = mapper(comment)
	}
	return res
}
