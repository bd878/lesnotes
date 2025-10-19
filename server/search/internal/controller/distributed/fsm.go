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

type Repository interface {
	SaveMessage(ctx context.Context, id, userID int64, name, title, text string, private bool) error
	DeleteMessage(ctx context.Context, id, userID int64) error
	SearchMessages(ctx context.Context, userID int64, substr string) (list []*searchmodel.Message, err error)
	Truncate(ctx context.Context) error
	Dump(ctx context.Context) (reader io.ReadCloser, err error)
	Restore(ctx context.Context, reader io.ReadCloser) (err error)
}

var _ raft.FSM = (*fsm)(nil)

type fsm struct {
	repo Repository
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

	logger.Debugw("apply append", "id", cmd.Id)

	err := f.repo.SaveMessage(context.Background(), cmd.Id, cmd.UserId, cmd.Name, cmd.Title, cmd.Text, cmd.Private)
	if err != nil {
		return err
	}

	return nil
}

func (f *fsm) applyUpdate(raw []byte) interface{} {
	// not implemented
	return nil
}

func (f *fsm) applyDelete(raw []byte) interface{} {
	var cmd DeleteCommand
	proto.Unmarshal(raw, &cmd)

	logger.Debugw("apply delete", "id", cmd.Id)

	err := f.repo.DeleteMessage(context.Background(), cmd.Id, cmd.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (f *fsm) applyPublish(raw []byte) interface{} {
	// not implemented
	return nil
}

func (f *fsm) applyPrivate(raw []byte) interface{} {
	// not implemented
	return nil
}
