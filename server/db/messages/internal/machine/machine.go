package machine

import (
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"

	"github.com/bd878/gallery/server/api/messages"
	"github.com/bd878/gallery/server/api/comments"
	"github.com/bd878/gallery/server/api/translations"
	"log/slog"
	"github.com/bd878/gallery/server/db/messages/pkg/machine"
)

type MessagesRepository interface {
	Create(ctx context.Context, id int64, text, title string, userID int64, private bool, name, createdAt, updatedAt string) (err error)
	Update(ctx context.Context, userID, id int64, text, title, name *string, updatedAt string) (err error)
	DeleteMessage(ctx context.Context, userID, id int64) (err error)
	Publish(ctx context.Context, userID int64, ids []int64, updatedAt string) (err error)
	Private(ctx context.Context, userID int64, ids []int64, updatedAt string) (err error)
	DeleteUserMessages(ctx context.Context, userID int64) (err error)
}

type TranslationsRepository interface {
	SaveTranslation(ctx context.Context, messageID int64, lang, text, title, createdAt, updatedAt string) (err error)
	UpdateTranslation(ctx context.Context, messageID int64, lang string, text, title *string, updatedAt string) (err error)
	DeleteTranslation(ctx context.Context, messageID int64, lang string) (err error)
	DeleteMessage(ctx context.Context, messageID int64) (err error)
}

type CommentsRepository interface {
	Create(ctx context.Context, id, userID, messageID int64, text string, metadata []byte, createdAt, updatedAt string) (err error)
	Update(ctx context.Context, id, userID int64, text *string, updatedAt string) (err error)
	Delete(ctx context.Context, id, userID int64) (err error)
	DeleteMessageComments(ctx context.Context, messageID int64) (err error)
}

type Dumper interface {
	Open(ctx context.Context) (ch chan *messages.MessagesSnapshot, err error)
	Restore(ctx context.Context, user *messages.MessagesSnapshot) (err error)
	Close() (err error)
}

var _ raft.FSM = (*Machine)(nil)

type Machine struct {
	dumper           Dumper
	messagesRepo     MessagesRepository
	commentsRepo     CommentsRepository
	translationsRepo TranslationsRepository
}

func New(messagesRepo MessagesRepository, translationsRepo TranslationsRepository,
	commentsRepo CommentsRepository, dumper Dumper) *Machine {
	return &Machine{
		dumper:           dumper,
		commentsRepo:     commentsRepo,
		messagesRepo:     messagesRepo,
		translationsRepo: translationsRepo,
	}
}

func (f *Machine) Apply(record *raft.Log) interface{} {
	buf := record.Data
	reqType := machine.RequestType(buf[0])
	switch reqType {
	case machine.AppendRequest:
		return f.applyAppend(buf[1:])
	case machine.UpdateRequest:
		return f.applyUpdate(buf[1:])
	case machine.DeleteUserMessagesRequest:
		return f.applyDeleteUserMessages(buf[1:])
	case machine.DeleteRequest:
		return f.applyDelete(buf[1:])
	case machine.PublishRequest:
		return f.applyPublish(buf[1:])
	case machine.PrivateRequest:
		return f.applyPrivate(buf[1:])
	case machine.AppendTranslationRequest:
		return f.applyAppendTranslation(buf[1:])
	case machine.UpdateTranslationRequest:
		return f.applyUpdateTranslation(buf[1:])
	case machine.DeleteTranslationRequest:
		return f.applyDeleteTranslation(buf[1:])
	case machine.AppendCommentRequest:
		return f.applyAppendComment(buf[1:])
	case machine.UpdateCommentRequest:
		return f.applyUpdateComment(buf[1:])
	case machine.DeleteCommentRequest:
		return f.applyDeleteComment(buf[1:])
	case machine.DeleteMessageCommentsRequest:
		return f.applyDeleteMessageComments(buf[1:])
	default:
		slog.Error("unknown request type", slog.Any("type", reqType))
	}
	return nil
}

func (f *Machine) applyAppend(raw []byte) interface{} {
	var cmd messages.AppendCommand
	proto.Unmarshal(raw, &cmd)

	// Put does not put message with same id twice
	return f.messagesRepo.Create(context.TODO(), cmd.Id, cmd.Text, cmd.Title, cmd.UserId, cmd.Private, cmd.Name, cmd.CreatedAt, cmd.UpdatedAt)
}

func (f *Machine) applyUpdate(raw []byte) interface{} {
	var cmd messages.UpdateCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.Update(context.TODO(), cmd.UserId, cmd.Id, cmd.Text, cmd.Title, cmd.Name, cmd.UpdatedAt)
}

func (f *Machine) applyDeleteUserMessages(raw []byte) interface{} {
	var cmd messages.DeleteUserMessagesCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.DeleteUserMessages(context.TODO(), cmd.UserId)
}

func (f *Machine) applyDelete(raw []byte) interface{} {
	var cmd messages.DeleteCommand
	proto.Unmarshal(raw, &cmd)

	err := f.messagesRepo.DeleteMessage(context.TODO(), cmd.UserId, cmd.Id)
	if err != nil {
		return err
	}

	err = f.translationsRepo.DeleteMessage(context.TODO(), cmd.Id)
	if err != nil {
		return err
	}

	return nil
}

func (f *Machine) applyPublish(raw []byte) interface{} {
	var cmd messages.PublishCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.Publish(context.TODO(), cmd.UserId, cmd.Ids, cmd.UpdatedAt)
}

func (f *Machine) applyPrivate(raw []byte) interface{} {
	var cmd messages.PrivateCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.Private(context.TODO(), cmd.UserId, cmd.Ids, cmd.UpdatedAt)
}

func (f *Machine) applyAppendTranslation(raw []byte) interface{} {
	var cmd translations.AppendTranslationCommand
	proto.Unmarshal(raw, &cmd)

	return f.translationsRepo.SaveTranslation(context.TODO(), cmd.MessageId, cmd.Lang, cmd.Title, cmd.Text, cmd.CreatedAt, cmd.UpdatedAt)
}

func (f *Machine) applyUpdateTranslation(raw []byte) interface{} {
	var cmd translations.UpdateTranslationCommand
	proto.Unmarshal(raw, &cmd)

	return f.translationsRepo.UpdateTranslation(context.TODO(), cmd.MessageId, cmd.Lang, cmd.Title, cmd.Text, cmd.UpdatedAt)
}

func (f *Machine) applyDeleteTranslation(raw []byte) interface{} {
	var cmd translations.DeleteTranslationCommand
	proto.Unmarshal(raw, &cmd)

	return f.translationsRepo.DeleteTranslation(context.TODO(), cmd.MessageId, cmd.Lang)
}

func (f *Machine) applyAppendComment(raw []byte) interface{} {
	var cmd comments.AppendCommentCommand
	proto.Unmarshal(raw, &cmd)

	return f.commentsRepo.Create(context.TODO(), cmd.Id, cmd.UserId, cmd.MessageId, cmd.Text, cmd.Metadata, cmd.CreatedAt, cmd.UpdatedAt)
}

func (f *Machine) applyUpdateComment(raw []byte) interface{} {
	var cmd comments.UpdateCommentCommand
	proto.Unmarshal(raw, &cmd)

	return f.commentsRepo.Update(context.TODO(), cmd.Id, cmd.UserId, cmd.Text, cmd.UpdatedAt)
}

func (f *Machine) applyDeleteComment(raw []byte) interface{} {
	var cmd comments.DeleteCommentCommand
	proto.Unmarshal(raw, &cmd)

	return f.commentsRepo.Delete(context.TODO(), cmd.Id, cmd.UserId)
}

func (f *Machine) applyDeleteMessageComments(raw []byte) interface{} {
	var cmd comments.DeleteMessageCommentsCommand
	proto.Unmarshal(raw, &cmd)

	return f.commentsRepo.DeleteMessageComments(context.TODO(), cmd.MessageId)
}
