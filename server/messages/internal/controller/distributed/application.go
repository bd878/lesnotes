package application

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/messages/internal/machine"
	"github.com/bd878/gallery/server/messages/internal/controller"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

type (
	MessagesRepository interface {
		Read(ctx context.Context, userIDs []int64, id int64, name string) (message *api.Message, err error)
		ReadMessages(ctx context.Context, userID int64, limit, offset int32) (messages []*api.Message, isLastPage bool, err error)
		ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*api.Message, err error)
	}

	TranslationsRepository interface {
		ReadTranslation(ctx context.Context, messageID int64, lang string) (translation *api.Translation, err error)
		ReadMessageTranslations(ctx context.Context, messageID int64) (translations []*api.TranslationPreview, err error)
		ListTranslations(ctx context.Context, messageID int64) (translations []*api.Translation, err error)
	}

	CommentsRepository interface {
		Read(ctx context.Context, id, userID int64) (comment *api.Comment, err error)
		ListMessageComments(ctx context.Context, messageID int64, limit, offset int32) (list *api.CommentsList, err error)
	}

	FilesGateway interface {
		ReadMessageFiles(ctx context.Context, id int64, userIDs []int64) (list []*api.File, err error)
	}

	Consensus interface {
		Apply(cmd []byte, timeout time.Duration) (err error)
		GetServers(ctx context.Context) ([]*api.Server, error)
	}

	App interface {
		ReadMessage(ctx context.Context, id int64, name string, userIDs []int64) (message *api.Message, err error)
		ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (messages []*api.Message, isLastPage bool, err error)
		ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*api.Message, err error)
		ReadTranslation(ctx context.Context, userID, messageID int64, lang string, name string) (translation *api.Translation, err error)
		ListTranslations(ctx context.Context, userID, messageID int64, name string) (translations []*api.Translation, err error)
		ReadComment(ctx context.Context, id, userID int64) (comment *api.Comment, err error)
		ListComments(ctx context.Context, userID, messageID *int64, name *string, limit, offset int32) (comments *api.CommentsList, err error)
		GetServers(ctx context.Context) ([]*api.Server, error)
	}

	Distributed struct {
		consensus        Consensus
		log              *logger.Logger
		commentsRepo     CommentsRepository
		messagesRepo     MessagesRepository
		translationsRepo TranslationsRepository
		filesGateway     FilesGateway
	}
)

var _ App = (*Distributed)(nil)

func New(consensus Consensus, messagesRepo MessagesRepository,
	translationsRepo TranslationsRepository, commentsRepo CommentsRepository, filesGateway FilesGateway, log *logger.Logger) *Distributed {
	return &Distributed{
		log:              log,
		consensus:        consensus,
		commentsRepo:     commentsRepo,
		messagesRepo:     messagesRepo,
		filesGateway:     filesGateway,
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
func (m *Distributed) ReadMessage(ctx context.Context, id int64, name string, userIDs []int64) (message *api.Message, err error) {
	m.log.Debugw("read message", "id", id, "name", name, "user_ids", userIDs)

	if id == 0 && name == "" {
		return nil, controller.ErrMessageIsRoot
	}

	message, err = m.messagesRepo.Read(ctx, userIDs, id, name)
	if err != nil {
		return
	}

	message.Files, err = m.filesGateway.ReadMessageFiles(ctx, message.Id, append(userIDs, message.UserId))
	if err != nil {
		return
	}

	message.Translations, err = m.translationsRepo.ReadMessageTranslations(ctx, message.Id)
	if err != nil {
		return
	}

	return
}

func (m *Distributed) ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (messages []*api.Message, isLastPage bool, err error) {
	m.log.Debugw("read messages", "user_id", userID, "limit", limit, "offset", offset, "ascending", ascending)

	messages, isLastPage, err = m.messagesRepo.ReadMessages(ctx, userID, limit, offset)
	if err != nil {
		return
	}

	// TODO: do not load files on each message in a list
	for _, message := range messages {
		message.Files, err = m.filesGateway.ReadMessageFiles(ctx, message.Id, []int64{userID, message.UserId})
		if err != nil {
			return
		}

		message.Translations, err = m.translationsRepo.ReadMessageTranslations(ctx, message.Id)
		if err != nil {
			return
		}
	}

	return
}

func (m *Distributed) ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*api.Message, err error) {
	m.log.Debugw("read batch messages", "user_id", userID, "ids", ids)

	messages, err = m.messagesRepo.ReadBatchMessages(ctx, userID, ids)
	if err != nil {
		return
	}

	for _, message := range messages {
		message.Files, err = m.filesGateway.ReadMessageFiles(ctx, message.Id, []int64{userID, message.UserId})
		if err != nil {
			return
		}

		message.Translations, err = m.translationsRepo.ReadMessageTranslations(ctx, message.Id)
		if err != nil {
			return
		}
	}

	return
}

func (m *Distributed) ReadTranslation(ctx context.Context, userID, messageID int64, lang string, name string) (translation *api.Translation, err error) {
	m.log.Debugw("read translation", "user_id", userID, "message_id", messageID, "lang", lang, "name", name)

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

func (m *Distributed) ListTranslations(ctx context.Context, userID, messageID int64, name string) (translations []*api.Translation, err error) {
	m.log.Debugw("list translations", "user_id", userID, "message_id", messageID, "name", name)

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

func (m *Distributed) ReadComment(ctx context.Context, id, userID int64) (comment *api.Comment, err error) {
	m.log.Debugw("read commnet", "id", id, "user_id", userID)
	return m.commentsRepo.Read(ctx, id, userID)
}

func (m *Distributed) ListComments(ctx context.Context, userID, messageID *int64, name *string, limit, offset int32) (comments *api.CommentsList, err error) {
	m.log.Debugw("list comments", "message_id", messageID, "user_id", userID, "name", name, "limit", limit, "offset", offset)

	if messageID != nil {
		return m.commentsRepo.ListMessageComments(ctx, *messageID, limit, offset)
	} else if name != nil {

		var message *api.Message

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
	m.log.Debugln("get servers")
	return m.consensus.GetServers(ctx)
}
