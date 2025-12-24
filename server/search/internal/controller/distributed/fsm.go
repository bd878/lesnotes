package distributed

import (
	"io"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	searchmodel "github.com/bd878/gallery/server/search/pkg/model"
	"github.com/bd878/gallery/server/logger"
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
	SearchMessages(ctx context.Context, userID int64, substr string, threadID int64, public int) (list []*searchmodel.Message, err error)
	Truncate(ctx context.Context) (err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

type FilesRepository interface {
	Truncate(ctx context.Context) (err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

type ThreadsRepository interface {
	Truncate(ctx context.Context) (err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

var _ raft.FSM = (*fsm)(nil)

type fsm struct {
	messagesRepo  MessagesRepository
	filesRepo     FilesRepository
	threadsRepo   ThreadsRepository
}

func (f *fsm) Apply(record *raft.Log) interface{} {
	buf := record.Data
	reqType := RequestType(buf[0])
	switch reqType {
	case AppendRequest:
		return f.applyAppend(buf[1:])
	case UpdateRequest:
		return f.applyUpdate(buf[1:])
	case DeleteRequest:
		return f.applyDelete(buf[1:])
	case PublishRequest:
		return f.applyPublish(buf[1:])
	case PrivateRequest:
		return f.applyPrivate(buf[1:])
	default:
		logger.Errorw("unknown request type", "type", reqType)
	}
	return nil
}

func (f *fsm) applyAppend(raw []byte) interface{} {
	var cmd AppendCommand
	proto.Unmarshal(raw, &cmd)

	err := f.messagesRepo.SaveMessage(context.Background(), cmd.Id, cmd.UserId, cmd.Name, cmd.Title, cmd.Text, cmd.Private)
	if err != nil {
		return err
	}

	return nil
}

func (f *fsm) applyUpdate(raw []byte) interface{} {
	var cmd UpdateCommand
	proto.Unmarshal(raw, &cmd)

	err := f.messagesRepo.UpdateMessage(context.Background(), cmd.Id, cmd.UserId, cmd.Name, cmd.Title, cmd.Text)
	if err != nil {
		return err
	}

	return nil
}

func (f *fsm) applyDelete(raw []byte) interface{} {
	var cmd DeleteCommand
	proto.Unmarshal(raw, &cmd)

	err := f.messagesRepo.DeleteMessage(context.Background(), cmd.Id, cmd.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (f *fsm) applyPublish(raw []byte) interface{} {
	var cmd PublishCommand
	proto.Unmarshal(raw, &cmd)

	err := f.messagesRepo.PublishMessages(context.Background(), cmd.Ids, cmd.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (f *fsm) applyPrivate(raw []byte) interface{} {
	var cmd PrivateCommand
	proto.Unmarshal(raw, &cmd)

	err := f.messagesRepo.PrivateMessages(context.Background(), cmd.Ids, cmd.UserId)
	if err != nil {
		return err
	}

	return nil
}
