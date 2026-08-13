package service

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/translations"
	"github.com/bd878/gallery/server/db/messages/pkg/machine"
	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/messages/internal/domain"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

type TranslationsController struct {
	client    translations.TranslationsClient
	publisher ddd.EventPublisher[ddd.Event]
}

func NewTranslationsController(container di.Container, publisher ddd.EventPublisher[ddd.Event]) *TranslationsController {
	client := container.Get("translationsClient").(translations.TranslationsClient)

	return &TranslationsController{
		client: client,
		publisher: publisher,
	}
}

func (s *TranslationsController) SaveTranslation(ctx context.Context, userID, messageID int64, lang, title, text string) (err error) {
	slog.Debug("save translation", slog.Int64("user_id", userID), slog.Int64("message_id", messageID), slog.String("lang", lang), slog.String("title", title), slog.String("text", text))

	createdAt := time.Now().UTC().Format(time.RFC3339)
	updatedAt := time.Now().UTC().Format(time.RFC3339)

	event, err := domain.CreateTranslation(userID, messageID, lang, title, text, createdAt, updatedAt)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return
	}

	cmd, err := proto.Marshal(&translations.AppendTranslationCommand{
		MessageId: messageID,
		Lang:      lang,
		Title:     title,
		Text:      text,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.AppendTranslationRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}

func (s *TranslationsController) UpdateTranslation(ctx context.Context, messageID int64, lang string, title, text *string) (err error) {

	logValues := []any{slog.Int64("message_id", messageID), slog.String("lang", lang)}
	if title != nil {
		logValues = append(logValues, slog.String("title", *title))
	}
	if text != nil {
		logValues = append(logValues, slog.String("text", *text))
	}
	slog.Debug("update translation", logValues...)

	updatedAt := time.Now().UTC().Format(time.RFC3339)

	event, err := domain.UpdateTranslation(messageID, lang, title, text, updatedAt)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return
	}

	cmd, err := proto.Marshal(&translations.UpdateTranslationCommand{
		MessageId: messageID,
		Lang:      lang,
		Title:     title,
		Text:      text,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.UpdateTranslationRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}

func (s *TranslationsController) DeleteTranslation(ctx context.Context, messageID int64, lang string) (err error) {
	slog.Debug("delete translation", slog.Int64("message_id", messageID), slog.String("lang", lang))

	event, err := domain.DeleteTranslation(messageID, lang)
	if err != nil {
		return err
	}

	err = s.publisher.Publish(ctx, event)
	if err != nil {
		return
	}

	cmd, err := proto.Marshal(&translations.DeleteTranslationCommand{
		MessageId: messageID,
		Lang:      lang,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType:  int32(machine.DeleteTranslationRequest),
		Cmd:      cmd,
		Duration: "10s",
	})

	return
}

func (s *TranslationsController) ReadTranslation(ctx context.Context, userID, messageID int64, lang string, name *string) (translation *model.Translation, err error) {

	logValues := []any{slog.Int64("user_id", userID), slog.Int64("message_id", messageID), slog.String("lang", lang)}
	if name != nil {
		logValues = append(logValues, slog.String("name", *name))
	}
	slog.Debug("read translation", logValues...)

	resp, err := s.client.ReadTranslation(ctx, &translations.ReadTranslationRequest{
		UserId: userID,
		Id:     messageID,
		Lang:   lang,
		Name:   name,
	})
	if err != nil {
		return nil, err
	}

	translation = model.TranslationFromProto(resp.Translation)

	return
}

func (s *TranslationsController) ListTranslations(ctx context.Context, userID, messageID int64, name string) (list []*model.Translation, err error) {
	slog.Debug("list translations", slog.Int64("user_id", userID), slog.Int64("message_id", messageID), slog.String("name", name))

	resp, err := s.client.ListTranslations(ctx, &translations.ListTranslationsRequest{
		UserId: userID,
		Id:     messageID,
		Name:   name,
	})
	if err != nil {
		return nil, err
	}

	list = model.MapTranslationsFromProto(model.TranslationFromProto, resp.Translations)

	return
}
