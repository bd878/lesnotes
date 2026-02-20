package application

import (
	"time"
	"errors"
	"context"
	"bytes"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/messages/internal/machine"
	"github.com/bd878/gallery/server/messages/internal/domain"
)

type MessagesRepository interface {
	Read(ctx context.Context, userIDs []int64, id int64, name string) (message *model.Message, err error)
	ReadMessages(ctx context.Context, userID int64, limit, offset int32) (messages []*model.Message, isLastPage bool, err error)
	ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*model.Message, err error)
}

type FilesRepository interface {
	ReadMessageFiles(ctx context.Context, messageID int64, userIDs []int64) (fileIDs []int64, err error)
}

type TranslationsRepository interface {
	ReadTranslation(ctx context.Context, messageID int64, lang string) (translation *model.Translation, err error)
	ReadMessageTranslations(ctx context.Context, messageID int64) (translations []*model.TranslationPreview, err error)
	ListTranslations(ctx context.Context, messageID int64) (translations []*model.Translation, err error)
}

type Consensus interface {
	Apply(cmd []byte, timeout time.Duration) (err error)
	GetServers(ctx context.Context) ([]*api.Server, error)
}

type Distributed struct {
	consensus         Consensus
	log               *logger.Logger
	publisher         ddd.EventPublisher[ddd.Event]
	messagesRepo      MessagesRepository
	filesRepo         FilesRepository
	translationsRepo  TranslationsRepository
}

func New(consensus Consensus, publisher ddd.EventPublisher[ddd.Event], messagesRepo MessagesRepository,
	filesRepo FilesRepository, translationsRepo TranslationsRepository, log *logger.Logger) *Distributed {
	return &Distributed{
		log:              log,
		publisher:        publisher,
		consensus:        consensus,
		messagesRepo:     messagesRepo,
		filesRepo:        filesRepo,
		translationsRepo: translationsRepo,
	}
}

func (m *Distributed) apply(ctx context.Context, reqType machine.RequestType, cmd []byte) (err error) {
	var buf bytes.Buffer
	_, err = buf.Write([]byte{byte(reqType)})
	if err != nil {
		return
	}

	_, err = buf.Write(cmd)
	if err != nil {
		return
	}

	return m.consensus.Apply(buf.Bytes(), 10*time.Second)
}

func (m *Distributed) SaveMessage(ctx context.Context, id int64, text, title string, fileIDs []int64, userID int64, private bool, name string) (err error) {
	m.log.Debugw("save message", "id", id, "text", text, "title", title, "file_ids", fileIDs, "user_id", userID, "private", private, "name", name)

	event, err := domain.CreateMessage(id, text, title, fileIDs, userID, private, name)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&machine.AppendCommand{
		Id:       id,
		Text:     text,
		Title:    title,
		FileIds:  fileIDs,
		UserId:   userID,
		Private:  private,
		Name:     name,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.AppendRequest, cmd)
	if err != nil {
		return
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *Distributed) UpdateMessage(ctx context.Context, id int64, text, title, name string, fileIDs []int64, userID int64) (err error) {
	m.log.Debugw("update message", "id", id, "text", text, "title", title, "name", name, "file_ids", fileIDs, "user_id", userID)

	event, err := domain.UpdateMessage(id, text, title, fileIDs, userID, name)
	if err != nil {
		return err
	}

	cmd, err := proto.Marshal(&machine.UpdateCommand{
		Id:       id,
		UserId:   userID,
		FileIds:  fileIDs,
		Text:     text,
		Name:     name,
		Title:    title,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.UpdateRequest, cmd)
	if err != nil {
		return err
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *Distributed) DeleteUserMessages(ctx context.Context, userID int64) (err error) {
	m.log.Debugw("delete user messages", "user_id", userID)

	cmd, err := proto.Marshal(&machine.DeleteUserMessagesCommand{
		UserId: userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.DeleteUserMessagesRequest, cmd)

	return
}

func (m *Distributed) DeleteFile(ctx context.Context, id, userID int64) (err error) {
	m.log.Debugw("delete file", "id", id, "user_id", userID)

	cmd, err := proto.Marshal(&machine.DeleteFileCommand{
		Id:      id,
		UserId:  userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.DeleteFileRequest, cmd)

	return
}

func (m *Distributed) DeleteMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	m.log.Debugw("delete messages", "ids", ids, "user_id", userID)

	for _, id := range ids {

		event, err := domain.DeleteMessage(id, userID)
		if err != nil {
			return err
		}

		cmd, err := proto.Marshal(&machine.DeleteCommand{
			Id:     id,
			UserId: userID,
		})
		if err != nil {
			return err
		}

		err = m.apply(ctx, machine.DeleteRequest, cmd)
		if err != nil {
			return err
		}

		m.publisher.Publish(context.Background(), event)

	}

	return
}

func (m *Distributed) PublishMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	m.log.Debugw("publish messages", "ids", ids, "user_id", userID)

	cmd, err := proto.Marshal(&machine.PublishCommand{
		Ids:           ids,
		UserId:        userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.PublishRequest, cmd)
	if err != nil {
		return
	}

	event, err := domain.PublishMessages(userID, ids)
	if err != nil {
		return err
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *Distributed) PrivateMessages(ctx context.Context, ids []int64, userID int64) (err error) {
	m.log.Debugw("private messages", "ids", ids, "user_id", userID)

	cmd, err := proto.Marshal(&machine.PrivateCommand{
		Ids:           ids,
		UserId:        userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.PrivateRequest, cmd)
	if err != nil {
		return
	}

	event, err := domain.PrivateMessages(userID, ids)
	if err != nil {
		return err
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *Distributed) SaveTranslation(ctx context.Context, userID, messageID int64, lang, title, text string) (err error) {
	m.log.Debugw("save translation", "user_id", userID, "message_id", messageID, "lang", lang, "title", title, "text", text)

	cmd, err := proto.Marshal(&machine.AppendTranslationCommand{
		MessageId:    messageID,
		Lang:         lang,
		Title:        title,
		Text:         text,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.AppendTranslationRequest, cmd)
	if err != nil {
		return
	}

	event, err := domain.CreateTranslation(userID, messageID, lang, title, text)
	if err != nil {
		return err
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *Distributed) UpdateTranslation(ctx context.Context, messageID int64, lang string, title, text *string) (err error) {
	// TODO: deal with nil pointer
	m.log.Debugw("update translation", "message_id", messageID, "lang", lang, "title", title, "text", text)

	cmd, err := proto.Marshal(&machine.UpdateTranslationCommand{
		MessageId:     messageID,
		Lang:          lang,
		Title:         title,
		Text:          text,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.UpdateTranslationRequest, cmd)
	if err != nil {
		return
	}

	event, err := domain.UpdateTranslation(messageID, lang, title, text)
	if err != nil {
		return err
	}

	return m.publisher.Publish(context.Background(), event)
}

func (m *Distributed) DeleteTranslation(ctx context.Context, messageID int64, lang string) (err error) {
	m.log.Debugw("delete translation", "message_id", messageID, "lang", lang)

	cmd, err := proto.Marshal(&machine.DeleteTranslationCommand{
		MessageId:  messageID,
		Lang:       lang,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.DeleteTranslationRequest, cmd)
	if err != nil {
		return
	}

	event, err := domain.DeleteTranslation(messageID, lang)
	if err != nil {
		return err
	}

	return m.publisher.Publish(context.Background(), event)
}

// TODO: pass one userID only, for public messages create ReadPublicMessage request
func (m *Distributed) ReadMessage(ctx context.Context, id int64, name string, userIDs []int64) (message *model.Message, err error) {
	m.log.Debugw("read message", "id", id, "name", name, "user_ids", userIDs)

	message, err = m.messagesRepo.Read(ctx, userIDs, id, name)
	if err != nil {
		return
	}

	message.FileIDs, err = m.filesRepo.ReadMessageFiles(ctx, message.ID /* cannot read by name */, append(userIDs, message.UserID))
	if err != nil {
		return
	}

	message.Translations, err = m.translationsRepo.ReadMessageTranslations(ctx, message.ID)
	if err != nil {
		return
	}

	return
}

func (m *Distributed) ReadMessages(ctx context.Context, userID int64, limit, offset int32, ascending bool) (messages []*model.Message, isLastPage bool, err error) {
	m.log.Debugw("read messages", "user_id", userID, "limit", limit, "offset", offset, "ascending", ascending)

	messages, isLastPage, err = m.messagesRepo.ReadMessages(ctx, userID, limit, offset)
	if err != nil {
		return
	}

	for _, message := range messages {
		message.FileIDs, err = m.filesRepo.ReadMessageFiles(ctx, message.ID, []int64{userID, message.UserID})
		if err != nil {
			return
		}

		message.Translations, err = m.translationsRepo.ReadMessageTranslations(ctx, message.ID)
		if err != nil {
			return
		}
	}

	return
}

func (m *Distributed) ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*model.Message, err error) {
	m.log.Debugw("read batch messages", "user_id", userID, "ids", ids)

	messages, err = m.messagesRepo.ReadBatchMessages(ctx, userID, ids)
	if err != nil {
		return
	}

	for _, message := range messages {
		message.FileIDs, err = m.filesRepo.ReadMessageFiles(ctx, message.ID, []int64{userID, message.UserID})
		if err != nil {
			return
		}

		message.Translations, err = m.translationsRepo.ReadMessageTranslations(ctx, message.ID)
		if err != nil {
			return
		}
	}

	return
}

func (m *Distributed) ReadTranslation(ctx context.Context, userID, messageID int64, lang string, name string) (translation *model.Translation, err error) {
	m.log.Debugw("read translation", "user_id", userID, "message_id", messageID, "lang", lang, "name", name)

	message, err := m.messagesRepo.Read(ctx, []int64{userID}, messageID, name)
	if err != nil {
		return nil, err
	}

	// only owner can read translations of his private message
	if message.Private && (message.UserID != userID) {
		return nil, errors.New("cannot read private message")
	}

	translation, err = m.translationsRepo.ReadTranslation(ctx, message.ID, lang)
	if err != nil {
		return nil, err
	}

	return
}

func (m *Distributed) ListTranslations(ctx context.Context, userID, messageID int64, name string) (translations []*model.Translation, err error) {
	m.log.Debugw("list translations", "user_id", userID, "message_id", messageID, "name", name)

	message, err := m.messagesRepo.Read(ctx, []int64{userID}, messageID, name)
	if err != nil {
		return nil, err
	}

	// only owner can list translations of his private message
	if message.Private && (message.UserID != userID) {
		return nil, errors.New("cannot list private translations")
	}

	translations, err = m.translationsRepo.ListTranslations(ctx, message.ID)
	if err != nil {
		return nil, err
	}

	return	
}

func (m *Distributed) GetServers(ctx context.Context) ([]*api.Server, error) {
	m.log.Debugln("get servers")
	return m.consensus.GetServers(ctx)
}
