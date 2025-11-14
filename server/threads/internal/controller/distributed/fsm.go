package distributed

import (
	"io"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/logger"
)

type RepoConnection interface {
	Release()
}

type Repository interface {
	CreateThread(ctx context.Context, id, userID int64, parentID int64, name string, private bool) error
	UpdateThread(ctx context.Context, id, userID int64, parentID int64, name string, private int32) error
	PrivateThread(ctx context.Context, id, userID int64) error
	PublishThread(ctx context.Context, id int64, userID int64) error
	DeleteThread(ctx context.Context, id, userID int64) error
	ResolveThread(ctx context.Context, id, userID int64) (ids []int64, err error)
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

	err := f.repo.CreateThread(context.Background(), cmd.Id, cmd.UserId, cmd.ParentId, cmd.Name, cmd.Private)
	if err != nil {
		return err
	}

	return nil
}

func (f *fsm) applyUpdate(raw []byte) interface{} {
	var cmd UpdateCommand
	proto.Unmarshal(raw, &cmd)

	err := f.repo.UpdateThread(context.Background(), cmd.Id, cmd.UserId, cmd.ParentId, cmd.Name, cmd.Private)
	if err != nil {
		return err
	}

	return nil
}

func (f *fsm) applyDelete(raw []byte) interface{} {
	var cmd DeleteCommand
	proto.Unmarshal(raw, &cmd)

	err := f.repo.DeleteThread(context.Background(), cmd.Id, cmd.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (f *fsm) applyPublish(raw []byte) interface{} {
	var cmd PublishCommand
	proto.Unmarshal(raw, &cmd)

	err := f.repo.PublishThread(context.Background(), cmd.Id, cmd.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (f *fsm) applyPrivate(raw []byte) interface{} {
	var cmd PrivateCommand
	proto.Unmarshal(raw, &cmd)

	err := f.repo.PrivateThread(context.Background(), cmd.Id, cmd.UserId)
	if err != nil {
		return err
	}

	return nil
}
