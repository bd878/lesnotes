package machine

import (
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
)

type MessagesRepository interface {
	SaveMessage(ctx context.Context, id, userID int64, name, title, text string, private bool, createdAt, updatedAt string) error
	UpdateMessage(ctx context.Context, id, userID int64, name, title, text *string, updatedAt string) error
	PrivateMessages(ctx context.Context, ids []int64, userID int64, updatedAt string) error
	PublishMessages(ctx context.Context, ids []int64, userID int64, updatedAt string) error
	DeleteMessage(ctx context.Context, id, userID int64) error
}

type FilesRepository interface {
	SaveFile(ctx context.Context, id, userID int64, name, description, mime string, private bool, size int64, createdAt, updatedAt string) (err error)
	DeleteFile(ctx context.Context, id, userID int64) (err error)
	PublishFile(ctx context.Context, id, userID int64, updatedAt string) (err error)
	PrivateFile(ctx context.Context, id, userID int64, updatedAt string) (err error)
}

type ThreadsRepository interface {
	SaveThread(ctx context.Context, id, userID, parentID int64, name, description string, private bool, createdAt, updatedAt string) error
	UpdateThread(ctx context.Context, id, userID int64, name, description *string, updatedAt string) error
	DeleteThread(ctx context.Context, id, userID int64) error
	ChangeThreadParent(ctx context.Context, id, userID, parentID int64) error
	PublishThread(ctx context.Context, id, userID int64, updatedAt string) error
	PrivateThread(ctx context.Context, id, userID int64, updatedAt string) error
}

type TranslationsRepository interface {
	SaveTranslation(ctx context.Context, userID, messageID int64, lang, title, text string, createdAt, updatedAt string) error
	DeleteTranslation(ctx context.Context, messageID int64, lang string) error
	UpdateTranslation(ctx context.Context, messageID int64, lang string, title, text *string, updatedAt string) error
}

type Dumper interface {
	Open(ctx context.Context) (ch chan *api.SearchSnapshot, err error)
	Restore(ctx context.Context, user *api.SearchSnapshot) (err error)
	Close() (err error)
}

var _ raft.FSM = (*Machine)(nil)

type Machine struct {
	log                *logger.Logger
	dumper             Dumper
	messagesRepo       MessagesRepository
	filesRepo          FilesRepository
	threadsRepo        ThreadsRepository
	translationsRepo   TranslationsRepository
}

func New(messagesRepo MessagesRepository, filesRepo FilesRepository, threadsRepo ThreadsRepository,
	translationsRepo TranslationsRepository, dumper Dumper, log *logger.Logger) *Machine {
	return &Machine{
		log:                 log,
		dumper:              dumper,
		messagesRepo:        messagesRepo,
		filesRepo:           filesRepo,
		translationsRepo:    translationsRepo,
		threadsRepo:         threadsRepo,
	}
}

func (f *Machine) Apply(record *raft.Log) interface{} {
	buf := record.Data
	reqType := RequestType(buf[0])
	switch reqType {
	case AppendMessageRequest:
		return f.applyAppendMessage(buf[1:])
	case UpdateMessageRequest:
		return f.applyUpdateMessage(buf[1:])
	case DeleteMessageRequest:
		return f.applyDeleteMessage(buf[1:])
	case PublishMessagesRequest:
		return f.applyPublishMessages(buf[1:])
	case PrivateMessagesRequest:
		return f.applyPrivateMessages(buf[1:])
	case AppendThreadRequest:
		return f.applyAppendThread(buf[1:])
	case UpdateThreadRequest:
		return f.applyUpdateThread(buf[1:])
	case DeleteThreadRequest:
		return f.applyDeleteThread(buf[1:])
	case ChangeThreadParentRequest:
		return f.applyChangeThreadParent(buf[1:])
	case PublishThreadRequest:
		return f.applyPublishThread(buf[1:])
	case PrivateThreadRequest:
		return f.applyPrivateThread(buf[1:])
	case AppendFileRequest:
		return f.applyAppendFile(buf[1:])
	case DeleteFileRequest:
		return f.applyDeleteFile(buf[1:])
	case PublishFileRequest:
		return f.applyPublishFile(buf[1:])
	case PrivateFileRequest:
		return f.applyPrivateFile(buf[1:])
	case AppendTranslationRequest:
		return f.applyAppendTranslation(buf[1:])
	case DeleteTranslationRequest:
		return f.applyDeleteTranslation(buf[1:])
	case UpdateTranslationRequest:
		return f.applyUpdateTranslation(buf[1:])
	default:
		logger.Errorw("unknown request type", "type", reqType)
	}
	return nil
}

func (f *Machine) applyAppendMessage(raw []byte) interface{} {
	var cmd AppendMessageCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.SaveMessage(context.TODO(), cmd.Id, cmd.UserId, cmd.Name, cmd.Title, cmd.Text, cmd.Private, cmd.CreatedAt, cmd.UpdatedAt)
}

func (f *Machine) applyUpdateMessage(raw []byte) interface{} {
	var cmd UpdateMessageCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.UpdateMessage(context.TODO(), cmd.Id, cmd.UserId, cmd.Name, cmd.Title, cmd.Text, cmd.UpdatedAt)
}

func (f *Machine) applyDeleteMessage(raw []byte) interface{} {
	var cmd DeleteMessageCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.DeleteMessage(context.TODO(), cmd.Id, cmd.UserId)
}

func (f *Machine) applyPublishMessages(raw []byte) interface{} {
	var cmd PublishMessagesCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.PublishMessages(context.TODO(), cmd.Ids, cmd.UserId, cmd.UpdatedAt)
}

func (f *Machine) applyPrivateMessages(raw []byte) interface{} {
	var cmd PrivateMessagesCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.PrivateMessages(context.TODO(), cmd.Ids, cmd.UserId, cmd.UpdatedAt)
}

func (f *Machine) applyAppendThread(raw []byte) interface{} {
	var cmd AppendThreadCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.SaveThread(context.TODO(), cmd.Id, cmd.UserId, cmd.ParentId, cmd.Name, cmd.Description, cmd.Private, cmd.CreatedAt, cmd.UpdatedAt)
}

func (f *Machine) applyUpdateThread(raw []byte) interface{} {
	var cmd UpdateThreadCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.UpdateThread(context.TODO(), cmd.Id, cmd.UserId, cmd.Name, cmd.Description, cmd.UpdatedAt)
}

func (f *Machine) applyDeleteThread(raw []byte) interface{} {
	var cmd DeleteThreadCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.DeleteThread(context.TODO(), cmd.Id, cmd.UserId)
}

func (f *Machine) applyChangeThreadParent(raw []byte) interface{} {
	var cmd ChangeThreadParentCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.ChangeThreadParent(context.TODO(), cmd.Id, cmd.UserId, cmd.ParentId)
}

func (f *Machine) applyPublishThread(raw []byte) interface{} {
	var cmd PublishThreadCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.PublishThread(context.TODO(), cmd.Id, cmd.UserId, cmd.UpdatedAt)
}

func (f *Machine) applyPrivateThread(raw []byte) interface{} {
	var cmd PrivateThreadCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.PrivateThread(context.TODO(), cmd.Id, cmd.UserId, cmd.UpdatedAt)
}

func (f *Machine) applyAppendFile(raw []byte) interface{} {
	var cmd AppendFileCommand
	proto.Unmarshal(raw, &cmd)

	return f.filesRepo.SaveFile(context.TODO(), cmd.Id, cmd.UserId, cmd.Name, cmd.Description, cmd.Mime, cmd.Private, cmd.Size, cmd.CreatedAt, cmd.UpdatedAt)
}

func (f *Machine) applyDeleteFile(raw []byte) interface{} {
	var cmd DeleteFileCommand
	proto.Unmarshal(raw, &cmd)

	return f.filesRepo.DeleteFile(context.TODO(), cmd.Id, cmd.UserId)
}

func (f *Machine) applyPublishFile(raw []byte) interface{} {
	var cmd PublishFileCommand
	proto.Unmarshal(raw, &cmd)

	return f.filesRepo.PublishFile(context.TODO(), cmd.Id, cmd.UserId, cmd.UpdatedAt)
}

func (f *Machine) applyPrivateFile(raw []byte) interface{} {
	var cmd PrivateFileCommand
	proto.Unmarshal(raw, &cmd)

	return f.filesRepo.PrivateFile(context.TODO(), cmd.Id, cmd.UserId, cmd.UpdatedAt)
}

func (f *Machine) applyAppendTranslation(raw []byte) interface{} {
	var cmd AppendTranslationCommand
	proto.Unmarshal(raw, &cmd)

	return f.translationsRepo.SaveTranslation(context.TODO(), cmd.UserId, cmd.MessageId, cmd.Lang, cmd.Title, cmd.Text, cmd.CreatedAt, cmd.UpdatedAt)
}

func (f *Machine) applyDeleteTranslation(raw []byte) interface{} {
	var cmd DeleteTranslationCommand
	proto.Unmarshal(raw, &cmd)

	return f.translationsRepo.DeleteTranslation(context.TODO(), cmd.MessageId, cmd.Lang)
}

func (f *Machine) applyUpdateTranslation(raw []byte) interface{} {
	var cmd UpdateTranslationCommand
	proto.Unmarshal(raw, &cmd)

	return f.translationsRepo.UpdateTranslation(context.TODO(), cmd.MessageId, cmd.Lang, cmd.Title, cmd.Text, cmd.UpdatedAt)
}
