package application

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/comments"
	"github.com/bd878/gallery/server/api/translations"
	"github.com/bd878/gallery/server/api/messages"
	"github.com/bd878/gallery/server/db/messages/pkg/machine"
	"github.com/bd878/gallery/server/db/messages/internal/controller"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

type (
	MessagesRepository interface {
		Read(ctx context.Context, userIDs []int64, id int64, name string) (message *messages.Message, err error)
		ReadMessages(ctx context.Context, userID int64, limit, offset int32) (messages []*messages.Message, isLastPage bool, err error)
		ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*messages.Message, err error)
	}

	TranslationsRepository interface {
		ReadTranslation(ctx context.Context, messageID int64, lang string) (translation *translations.Translation, err error)
		ReadMessageTranslations(ctx context.Context, messageID int64) (translations []*translations.TranslationPreview, err error)
		ListTranslations(ctx context.Context, messageID int64) (translations []*translations.Translation, err error)
	}

	CommentsRepository interface {
		Read(ctx context.Context, id, userID int64) (comment *comments.Comment, err error)
		ListMessageComments(ctx context.Context, messageID int64, limit, offset int32) (list *comments.CommentsList, err error)
	}

	Consensus interface {
		Apply(cmd []byte, timeout time.Duration) (err error)
		GetServers(ctx context.Context) ([]*api.Server, error)
	}

	App interface {
		ReadMessage(ctx context.Context, id int64, name string, userIDs []int64) (message *messages.Message, err error)
		ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (messages []*messages.Message, isLastPage bool, err error)
		ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*messages.Message, err error)
		ReadTranslation(ctx context.Context, userID, messageID int64, lang string, name string) (translation *translations.Translation, err error)
		ListTranslations(ctx context.Context, userID, messageID int64, name string) (translations []*translations.Translation, err error)
		ReadComment(ctx context.Context, id, userID int64) (comment *comments.Comment, err error)
		ListComments(ctx context.Context, userID, messageID *int64, name *string, limit, offset int32) (comments *comments.CommentsList, err error)
		GetServers(ctx context.Context) ([]*api.Server, error)
	}

	Distributed struct {
		consensus        Consensus
		commentsRepo     CommentsRepository
		messagesRepo     MessagesRepository
		translationsRepo TranslationsRepository
	}
)

var _ App = (*Distributed)(nil)

func New(consensus Consensus, messagesRepo MessagesRepository,
	translationsRepo TranslationsRepository, commentsRepo CommentsRepository) *Distributed {
	return &Distributed{
		consensus:        consensus,
		commentsRepo:     commentsRepo,
		messagesRepo:     messagesRepo,
		translationsRepo: translationsRepo,
	}
}

func (m *Distributed) Apply(ctx context.Context, reqType machine.RequestType, cmd []byte, duration time.Duration) (err error) {
	var buf bytes.Buffer
	_, err = buf.Write([]byte{byte(reqType)})
	if err != nil {
		return
	}

	_, err = buf.Write(cmd)
	if err != nil {
		return
	}

	return m.consensus.Apply(buf.Bytes(), duration)
}

// TODO: pass one userID only, for public messages create ReadPublicMessage request
func (m *Distributed) ReadMessage(ctx context.Context, id int64, name string, userIDs []int64) (message *messages.Message, err error) {
	slog.Debug("read message", slog.Int64("id", id), slog.String("name", name), slog.String("user_ids", fmt.Sprintf("%v", userIDs)))

	if id == 0 && name == "" {
		return nil, controller.ErrMessageIsRoot
	}

	message, err = m.messagesRepo.Read(ctx, userIDs, id, name)
	if err != nil {
		return
	}

	message.Translations, err = m.translationsRepo.ReadMessageTranslations(ctx, message.Id)
	if err != nil {
		return
	}

	return
}

func (m *Distributed) ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (messages []*messages.Message, isLastPage bool, err error) {
	slog.Debug("read messages", slog.Int64("user_id", userID), slog.Int("limit", int(limit)), slog.Int("offset", int(offset)), slog.Bool("ascending", ascending))

	messages, isLastPage, err = m.messagesRepo.ReadMessages(ctx, userID, limit, offset)
	if err != nil {
		return
	}

	// TODO: do not load files on each message in a list
	for _, message := range messages {
		message.Translations, err = m.translationsRepo.ReadMessageTranslations(ctx, message.Id)
		if err != nil {
			return
		}
	}

	return
}

func (m *Distributed) ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*messages.Message, err error) {
	slog.Debug("read batch messages", slog.Int64("user_id", userID), slog.String("ids", fmt.Sprintf("%v", ids)))

	messages, err = m.messagesRepo.ReadBatchMessages(ctx, userID, ids)
	if err != nil {
		return
	}

	for _, message := range messages {
		message.Translations, err = m.translationsRepo.ReadMessageTranslations(ctx, message.Id)
		if err != nil {
			return
		}
	}

	return
}

func (m *Distributed) ReadTranslation(ctx context.Context, userID, messageID int64, lang string, name string) (translation *translations.Translation, err error) {
	slog.Debug("read translation", slog.Int64("user_id", userID), slog.Int64("message_id", messageID), slog.String("lang", lang), slog.String("name", name))

	message, err := m.messagesRepo.Read(ctx, []int64{userID}, messageID, name)
	if err != nil {
		return nil, err
	}

	// only owner can read translations of his private message
	if message.Private && (message.UserId != userID) {
		return nil, errors.New("cannot read private message")
	}

	translation, err = m.translationsRepo.ReadTranslation(ctx, message.Id, lang)
	if err != nil {
		return nil, err
	}

	return
}

func (m *Distributed) ListTranslations(ctx context.Context, userID, messageID int64, name string) (translations []*translations.Translation, err error) {
	slog.Debug("list translations", slog.Int64("user_id", userID), slog.Int64("message_id", messageID), slog.String("name", name))

	message, err := m.messagesRepo.Read(ctx, []int64{userID}, messageID, name)
	if err != nil {
		return nil, err
	}

	// only owner can list translations of his private message
	if message.Private && (message.UserId != userID) {
		return nil, errors.New("cannot list private translations")
	}

	translations, err = m.translationsRepo.ListTranslations(ctx, message.Id)
	if err != nil {
		return nil, err
	}

	return
}

func (m *Distributed) ReadComment(ctx context.Context, id, userID int64) (comment *comments.Comment, err error) {
	slog.Debug("read commnet", slog.Int64("id", id), slog.Int64("user_id", userID))
	return m.commentsRepo.Read(ctx, id, userID)
}

func (m *Distributed) ListComments(ctx context.Context, userID, messageID *int64, name *string, limit, offset int32) (result *comments.CommentsList, err error) {
	logValues := []any{slog.Int("limit", int(limit)), slog.Int("offset", int(offset))}
	if messageID != nil {
		logValues = append(logValues, slog.Int64("message_id", *messageID))
	}
	if userID != nil {
		logValues = append(logValues, slog.Int64("user_id", *userID))
	}
	if name != nil {
		logValues = append(logValues, slog.String("name", *name))
	}
	slog.Debug("list comments", logValues...)

	if messageID != nil {
		return m.commentsRepo.ListMessageComments(ctx, *messageID, limit, offset)
	} else if name != nil {

		var message *messages.Message

		if userID != nil {
			message, err = m.messagesRepo.Read(ctx, []int64{*userID}, 0, *name)
		} else {
			message, err = m.messagesRepo.Read(ctx, []int64{users.PublicUserID}, 0, *name)
		}
		if err != nil {
			return
		}

		return m.commentsRepo.ListMessageComments(ctx, message.Id, limit, offset)

	} else {
		return nil, errors.New("list users comments not implemented")
	}
}

func (m *Distributed) GetServers(ctx context.Context) ([]*api.Server, error) {
	slog.Debug("get servers")
	return m.consensus.GetServers(ctx)
}
