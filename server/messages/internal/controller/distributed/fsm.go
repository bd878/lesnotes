package distributed

import (
	"io"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/logger"
)

type RepoConnection interface {
	Release()
}

type MessagesRepository interface {
	Create(ctx context.Context, id int64, text, title string, userID int64, private bool, name string) (err error)
	Update(ctx context.Context, userID, id int64, newText, newTitle, newName string) (err error)
	DeleteMessage(ctx context.Context, userID, id int64) (err error)
	Publish(ctx context.Context, userID int64, ids []int64) (err error)
	Private(ctx context.Context, userID int64, ids []int64) (err error)
	Read(ctx context.Context, userIDs []int64, id int64, name string) (message *model.Message, err error)
	DeleteUserMessages(ctx context.Context, userID int64) (err error)
	ReadMessages(ctx context.Context, userID int64, limit, offset int32) (messages []*model.Message, isLastPage bool, err error)
	ReadBatchMessages(ctx context.Context, userID int64, ids []int64) (messages []*model.Message, err error)
	Truncate(ctx context.Context) error
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

type FilesRepository interface {
	DeleteFile(ctx context.Context, id, userID int64) (err error)
	SaveMessageFiles(ctx context.Context, messageID, userID int64, fileIDs []int64) (err error)
	UpdateMessageFiles(ctx context.Context, messageID, userID int64, fileIDs []int64) (err error)
	DeleteMessage(ctx context.Context, messageID, userID int64) (err error)
	Truncate(ctx context.Context) error
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

var _ raft.FSM = (*fsm)(nil)

type fsm struct {
	messagesRepo MessagesRepository
	filesRepo    FilesRepository
}

/**
 * Returns empty interface. It is either an error,
 * or new msg with unique id, saved in repo.
 * 
 * Apply replicates log state from the bottom up.
 * Leader makes Apply on start.
 */
func (f *fsm) Apply(record *raft.Log) interface{} {
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
	default:
		logger.Errorw("unknown request type", "type", reqType)
	}
	return nil
}

func (f *fsm) applyAppend(raw []byte) interface{} {
	var cmd AppendCommand
	proto.Unmarshal(raw, &cmd)

	// Put does not put message with same id twice
	err := f.messagesRepo.Create(context.Background(), cmd.Id, cmd.Text, cmd.Title, cmd.UserId, cmd.Private, cmd.Name)
	if err != nil {
		return err
	}

	err = f.filesRepo.SaveMessageFiles(context.Background(), cmd.Id, cmd.UserId, cmd.FileIds)
	if err != nil {
		return err
	}

	return nil
}

func (f *fsm) applyUpdate(raw []byte) interface{} {
	var cmd UpdateCommand
	proto.Unmarshal(raw, &cmd)

	err := f.messagesRepo.Update(context.Background(), cmd.UserId, cmd.Id, cmd.Text, cmd.Title, cmd.Name)
	if err != nil {
		return err
	}

	err = f.filesRepo.UpdateMessageFiles(context.Background(), cmd.Id, cmd.UserId, cmd.FileIds)
	if err != nil {
		return err
	}

	return nil
}

func (f *fsm) applyDeleteUserMessages(raw []byte) interface{} {
	var cmd DeleteUserMessagesCommand
	proto.Unmarshal(raw, &cmd)

	err := f.messagesRepo.DeleteUserMessages(context.Background(), cmd.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (f *fsm) applyDelete(raw []byte) interface{} {
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

	return nil
}

func (f *fsm) applyPublish(raw []byte) interface{} {
	var cmd PublishCommand
	proto.Unmarshal(raw, &cmd)

	err := f.messagesRepo.Publish(context.Background(), cmd.UserId, cmd.Ids)
	if err != nil {
		return err
	}

	return nil
}

func (f *fsm) applyPrivate(raw []byte) interface{} {
	var cmd PrivateCommand
	proto.Unmarshal(raw, &cmd)

	err := f.messagesRepo.Private(context.Background(), cmd.UserId, cmd.Ids)
	if err != nil {
		return err
	}

	return nil
}

func (f *fsm) applyDeleteFile(raw []byte) interface{} {
	var cmd DeleteFileCommand
	proto.Unmarshal(raw, &cmd)

	err := f.filesRepo.DeleteFile(context.Background(), cmd.Id, cmd.UserId)
	if err != nil {
		return err
	}

	return nil
}