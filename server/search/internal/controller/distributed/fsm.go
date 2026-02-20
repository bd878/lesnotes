package distributed

import (
	"io"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/search/pkg/model"
)

type RepoConnection interface {
	Release()
}

type MessagesRepository interface {
	SaveMessage(ctx context.Context, id, userID int64, name, title, text string, private bool) error
	UpdateMessage(ctx context.Context, id, userID int64, name, title, text string) error
	PrivateMessages(ctx context.Context, ids []int64, userID int64) error
	PublishMessages(ctx context.Context, ids []int64, userID int64) error
	DeleteMessage(ctx context.Context, id, userID int64) error
	SearchMessages(ctx context.Context, userID int64, substr string, public int) (list []*model.Message, err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

type FilesRepository interface {
	SaveFile(ctx context.Context, id, userID int64, name, description, mime string, private bool, size int64) (err error)
	DeleteFile(ctx context.Context, id, userID int64) (err error)
	PublishFile(ctx context.Context, id, userID int64) (err error)
	PrivateFile(ctx context.Context, id, userID int64) (err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

type ThreadsRepository interface {
	SaveThread(ctx context.Context, id, userID, parentID int64, name, description string, private bool) error
	UpdateThread(ctx context.Context, id, userID int64, name, description string) error
	DeleteThread(ctx context.Context, id, userID int64) error
	ChangeThreadParent(ctx context.Context, id, userID, parentID int64) error
	PublishThread(ctx context.Context, id, userID int64) error
	PrivateThread(ctx context.Context, id, userID int64) error
	SearchThreads(ctx context.Context, parentID, userID int64) (list []*model.Thread, err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

type TranslationsRepository interface {
	SaveTranslation(ctx context.Context, userID, messageID int64, lang, title, text string) error
	DeleteTranslation(ctx context.Context, messageID int64, lang string) error
	UpdateTranslation(ctx context.Context, messageID int64, lang string, title, text *string) error
	SearchTranslations(ctx context.Context, userID int64, substr string) (list []*model.Translation, err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

var _ raft.FSM = (*fsm)(nil)

type fsm struct {
	messagesRepo       MessagesRepository
	filesRepo          FilesRepository
	threadsRepo        ThreadsRepository
	translationsRepo   TranslationsRepository
}

func (f *fsm) Apply(record *raft.Log) interface{} {
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

func (f *fsm) applyAppendMessage(raw []byte) interface{} {
	var cmd AppendMessageCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.SaveMessage(context.Background(), cmd.Id, cmd.UserId, cmd.Name, cmd.Title, cmd.Text, cmd.Private)
}

func (f *fsm) applyUpdateMessage(raw []byte) interface{} {
	var cmd UpdateMessageCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.UpdateMessage(context.Background(), cmd.Id, cmd.UserId, cmd.Name, cmd.Title, cmd.Text)
}

func (f *fsm) applyDeleteMessage(raw []byte) interface{} {
	var cmd DeleteMessageCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.DeleteMessage(context.Background(), cmd.Id, cmd.UserId)
}

func (f *fsm) applyPublishMessages(raw []byte) interface{} {
	var cmd PublishMessagesCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.PublishMessages(context.Background(), cmd.Ids, cmd.UserId)
}

func (f *fsm) applyPrivateMessages(raw []byte) interface{} {
	var cmd PrivateMessagesCommand
	proto.Unmarshal(raw, &cmd)

	return f.messagesRepo.PrivateMessages(context.Background(), cmd.Ids, cmd.UserId)
}

func (f *fsm) applyAppendThread(raw []byte) interface{} {
	var cmd AppendThreadCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.SaveThread(context.Background(), cmd.Id, cmd.UserId, cmd.ParentId, cmd.Name, cmd.Description, cmd.Private)
}

func (f *fsm) applyUpdateThread(raw []byte) interface{} {
	var cmd UpdateThreadCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.UpdateThread(context.Background(), cmd.Id, cmd.UserId, cmd.Name, cmd.Description)
}

func (f *fsm) applyDeleteThread(raw []byte) interface{} {
	var cmd DeleteThreadCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.DeleteThread(context.Background(), cmd.Id, cmd.UserId)
}

func (f *fsm) applyChangeThreadParent(raw []byte) interface{} {
	var cmd ChangeThreadParentCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.ChangeThreadParent(context.Background(), cmd.Id, cmd.UserId, cmd.ParentId)
}

func (f *fsm) applyPublishThread(raw []byte) interface{} {
	var cmd PublishThreadCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.PublishThread(context.Background(), cmd.Id, cmd.UserId)
}

func (f *fsm) applyPrivateThread(raw []byte) interface{} {
	var cmd PrivateThreadCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.PrivateThread(context.Background(), cmd.Id, cmd.UserId)
}

func (f *fsm) applyAppendFile(raw []byte) interface{} {
	var cmd AppendFileCommand
	proto.Unmarshal(raw, &cmd)

	return f.filesRepo.SaveFile(context.Background(), cmd.Id, cmd.UserId, cmd.Name, cmd.Description, cmd.Mime, cmd.Private, cmd.Size)
}

func (f *fsm) applyDeleteFile(raw []byte) interface{} {
	var cmd DeleteFileCommand
	proto.Unmarshal(raw, &cmd)

	return f.filesRepo.DeleteFile(context.Background(), cmd.Id, cmd.UserId)
}

func (f *fsm) applyPublishFile(raw []byte) interface{} {
	var cmd PublishFileCommand
	proto.Unmarshal(raw, &cmd)

	return f.filesRepo.PublishFile(context.Background(), cmd.Id, cmd.UserId)
}

func (f *fsm) applyPrivateFile(raw []byte) interface{} {
	var cmd PrivateFileCommand
	proto.Unmarshal(raw, &cmd)

	return f.filesRepo.PrivateFile(context.Background(), cmd.Id, cmd.UserId)
}

func (f *fsm) applyAppendTranslation(raw []byte) interface{} {
	var cmd AppendTranslationCommand
	proto.Unmarshal(raw, &cmd)

	return f.translationsRepo.SaveTranslation(context.Background(), cmd.UserId, cmd.MessageId, cmd.Lang, cmd.Title, cmd.Text)
}

func (f *fsm) applyDeleteTranslation(raw []byte) interface{} {
	var cmd DeleteTranslationCommand
	proto.Unmarshal(raw, &cmd)

	return f.translationsRepo.DeleteTranslation(context.Background(), cmd.MessageId, cmd.Lang)
}

func (f *fsm) applyUpdateTranslation(raw []byte) interface{} {
	var cmd UpdateTranslationCommand
	proto.Unmarshal(raw, &cmd)

	return f.translationsRepo.UpdateTranslation(context.Background(), cmd.MessageId, cmd.Lang, cmd.Title, cmd.Text)
}
