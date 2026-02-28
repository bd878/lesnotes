package machine

import (
	"io"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/internal/logger"
)

type MessagesRepository interface {
	Create(ctx context.Context, id int64, text, title string, userID int64, private bool, name, createdAt, updatedAt string) (err error)
	Update(ctx context.Context, userID, id int64, text, title, name *string, updatedAt string) (err error)
	DeleteMessage(ctx context.Context, userID, id int64) (err error)
	Publish(ctx context.Context, userID int64, ids []int64, updatedAt string) (err error)
	Private(ctx context.Context, userID int64, ids []int64, updatedAt string) (err error)
	DeleteUserMessages(ctx context.Context, userID int64) (err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

type FilesRepository interface {
	DeleteFile(ctx context.Context, id, userID int64) (err error)
	SaveMessageFiles(ctx context.Context, messageID, userID int64, fileIDs []int64) (err error)
	UpdateMessageFiles(ctx context.Context, messageID, userID int64, fileIDs []int64) (err error)
	DeleteMessage(ctx context.Context, messageID, userID int64) (err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

type TranslationsRepository interface {
	SaveTranslation(ctx context.Context, messageID int64, lang, text, title, createdAt, updatedAt string) (err error)
	UpdateTranslation(ctx context.Context, messageID int64, lang string, text, title *string, updatedAt string) (err error)
	DeleteTranslation(ctx context.Context, messageID int64, lang string) (err error)
	DeleteMessage(ctx context.Context, messageID int64) (err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

var _ raft.FSM = (*Machine)(nil)

type Machine struct {
	log               *logger.Logger
	messagesRepo      MessagesRepository
	filesRepo         FilesRepository
	translationsRepo  TranslationsRepository
}

func New(messagesRepo MessagesRepository, filesRepo FilesRepository,
	translationsRepo TranslationsRepository, log *logger.Logger) *Machine {
	return &Machine{
		log:              log,
		messagesRepo:     messagesRepo,
		filesRepo:        filesRepo,
		translationsRepo: translationsRepo,
	}
}

func (f *Machine) Apply(record *raft.Log) interface{} {
	buf := record.Data
	reqType := RequestType(buf[0])
	switch reqType {
	case AppendRequest:
		return f.applyAppend(buf[1:])
	case UpdateRequest:
		return f.applyUpdate(buf[1:])
	case DeleteUserMessagesRequest:
		return f.applyDeleteUserMessages(buf[1:])
	case DeleteRequest:
		return f.applyDelete(buf[1:])
	case PublishRequest:
		return f.applyPublish(buf[1:])
	case PrivateRequest:
		return f.applyPrivate(buf[1:])
	case DeleteFileRequest:
		return f.applyDeleteFile(buf[1:])
	case AppendTranslationRequest:
		return f.applyAppendTranslation(buf[1:])
	case UpdateTranslationRequest:
		return f.applyUpdateTranslation(buf[1:])
	case DeleteTranslationRequest:
		return f.applyDeleteTranslation(buf[1:])
	default:
		f.log.Errorw("unknown request type", "type", reqType)
	}
	return nil
}

func (f *Machine) applyAppend(raw []byte) interface{} {
	var cmd AppendCommand
	proto.Unmarshal(raw, &cmd)

	// Put does not put message with same id twice
	err := f.messagesRepo.Create(context.Background(), cmd.Id, cmd.Text, cmd.Title, cmd.UserId, cmd.Private, cmd.Name, cmd.CreatedAt, cmd.UpdatedAt)
	if err != nil {
		return err
	}

	return f.filesRepo.SaveMessageFiles(context.Background(), cmd.Id, cmd.UserId, cmd.FileIds)
}

func (f *Machine) applyUpdate(raw []byte) interface{} {
	var cmd UpdateCommand
	proto.Unmarshal(raw, &cmd)

	err := f.messagesRepo.Update(context.Background(), cmd.UserId, cmd.Id, cmd.Text, cmd.Title, cmd.Name, cmd.UpdatedAt)
	if err != nil {
		return err
	}

	return f.filesRepo.UpdateMessageFiles(context.Background(), cmd.Id, cmd.UserId, cmd.FileIds)
}

func (f *Machine) applyDeleteUserMessages(raw []byte) interface{} {
	var cmd DeleteUserMessagesCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.DeleteUserMessages(context.Background(), cmd.UserId)
}

func (f *Machine) applyDelete(raw []byte) interface{} {
	var cmd DeleteCommand
	proto.Unmarshal(raw, &cmd)

	err := f.messagesRepo.DeleteMessage(context.Background(), cmd.UserId, cmd.Id)
	if err != nil {
		return err
	}

	err = f.filesRepo.DeleteMessage(context.Background(), cmd.Id, cmd.UserId)
	if err != nil {
		return err
	}

	err = f.translationsRepo.DeleteMessage(context.Background(), cmd.Id)
	if err != nil {
		return err
	}

	return nil
}

func (f *Machine) applyPublish(raw []byte) interface{} {
	var cmd PublishCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.Publish(context.Background(), cmd.UserId, cmd.Ids, cmd.UpdatedAt)
}

func (f *Machine) applyPrivate(raw []byte) interface{} {
	var cmd PrivateCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.Private(context.Background(), cmd.UserId, cmd.Ids, cmd.UpdatedAt)
}

func (f *Machine) applyDeleteFile(raw []byte) interface{} {
	var cmd DeleteFileCommand
	proto.Unmarshal(raw, &cmd)

	return f.filesRepo.DeleteFile(context.Background(), cmd.Id, cmd.UserId)
}

func (f *Machine) applyAppendTranslation(raw []byte) interface{} {
	var cmd AppendTranslationCommand
	proto.Unmarshal(raw, &cmd)

	return f.translationsRepo.SaveTranslation(context.Background(), cmd.MessageId, cmd.Lang, cmd.Title, cmd.Text, cmd.CreatedAt, cmd.UpdatedAt)
}

func (f *Machine) applyUpdateTranslation(raw []byte) interface{} {
	var cmd UpdateTranslationCommand
	proto.Unmarshal(raw, &cmd)

	return f.translationsRepo.UpdateTranslation(context.Background(), cmd.MessageId, cmd.Lang, cmd.Title, cmd.Text, cmd.UpdatedAt)
}

func (f *Machine) applyDeleteTranslation(raw []byte) interface{} {
	var cmd DeleteTranslationCommand
	proto.Unmarshal(raw, &cmd)

	err := f.translationsRepo.DeleteTranslation(context.Background(), cmd.MessageId, cmd.Lang)
	if err != nil {
		return err
	}

	return nil
}
