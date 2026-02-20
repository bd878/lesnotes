package machine

import (
	"io"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/threads/pkg/model"
)

type ThreadsRepository interface {
	ListThreads(ctx context.Context, userID, parentID int64, limit, offset int32, asc bool) (ids []*model.Thread, isLastPage bool, err error)
	ReadThread(ctx context.Context, id, userID int64, name string) (thread *model.Thread, err error)
	ResolveThread(ctx context.Context, id, userID int64) (ids []int64, err error)
	CountThreads(ctx context.Context, id, userID int64) (total int32, err error)
	AppendThread(ctx context.Context, id, userID, parentID, nextID, prevID int64, name, description string, private bool) (err error)
	UpdateThread(ctx context.Context, id, userID int64, name, description string) (err error)
	PrivateThread(ctx context.Context, id, userID int64) error
	PublishThread(ctx context.Context, id, userID int64) error
	DeleteThread(ctx context.Context, id, userID int64) error
	ReorderThread(ctx context.Context, id, userID, parentID, nextID, prevID int64) (err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

var _ raft.FSM = (*Machine)(nil)

type Machine struct {
	log               *logger.Logger
	threadsRepo       ThreadsRepository
}

func New(threadsRepo ThreadsRepository, log *logger.Logger) *Machine {
	return &Machine{
		log:           log,
		threadsRepo:   threadsRepo,
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
	case DeleteRequest:
		return f.applyDelete(buf[1:])
	case PublishRequest:
		return f.applyPublish(buf[1:])
	case PrivateRequest:
		return f.applyPrivate(buf[1:])
	case ReorderRequest:
		return f.applyReorder(buf[1:])
	default:
		f.log.Errorw("unknown request type", "type", reqType)
	}
	return nil
}

func (f *Machine) applyAppend(raw []byte) interface{} {
	var cmd AppendCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.AppendThread(context.Background(), cmd.Id, cmd.UserId, cmd.ParentId, cmd.NextId, cmd.PrevId, cmd.Name, cmd.Description, cmd.Private)
}

func (f *Machine) applyReorder(raw []byte) interface{} {
	var cmd ReorderCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.ReorderThread(context.Background(), cmd.Id, cmd.UserId, cmd.ParentId, cmd.NextId, cmd.PrevId)
}

func (f *Machine) applyUpdate(raw []byte) interface{} {
	var cmd UpdateCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.UpdateThread(context.Background(), cmd.Id, cmd.UserId, cmd.Name, cmd.Description)
}

func (f *Machine) applyDelete(raw []byte) interface{} {
	var cmd DeleteCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.DeleteThread(context.Background(), cmd.Id, cmd.UserId)
}

func (f *Machine) applyPublish(raw []byte) interface{} {
	var cmd PublishCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.PublishThread(context.Background(), cmd.Id, cmd.UserId)
}

func (f *Machine) applyPrivate(raw []byte) interface{} {
	var cmd PrivateCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.PrivateThread(context.Background(), cmd.Id, cmd.UserId)
}
