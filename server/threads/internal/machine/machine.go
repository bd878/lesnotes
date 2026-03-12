package machine

import (
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
)

type ThreadsRepository interface {
	AppendThread(ctx context.Context, id, userID, parentID, nextID, prevID int64, name, description string, private bool, createdAt, updatedAt string) (err error)
	UpdateThread(ctx context.Context, id, userID int64, name, description *string, updatedAt string) (err error)
	PrivateThread(ctx context.Context, id, userID int64, updatedAt string) error
	PublishThread(ctx context.Context, id, userID int64, updatedAt string) error
	DeleteThread(ctx context.Context, id, userID int64) error
	ReorderThread(ctx context.Context, id, userID, parentID, nextID, prevID int64, updatedAt string) (err error)
}

type Dumper interface {
	Open(ctx context.Context) (ch chan *api.ThreadsSnapshot, err error)
	Restore(ctx context.Context, thread *api.ThreadsSnapshot) (err error)
	Close() (err error)
}

var _ raft.FSM = (*Machine)(nil)

type Machine struct {
	log               *logger.Logger
	threadsRepo       ThreadsRepository
	dumper            Dumper
}

func New(threadsRepo ThreadsRepository, dumper Dumper, log *logger.Logger) *Machine {
	return &Machine{
		log:           log,
		dumper:        dumper,
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

	return f.threadsRepo.AppendThread(context.TODO(), cmd.Id, cmd.UserId, cmd.ParentId, cmd.NextId,
		cmd.PrevId, cmd.Name, cmd.Description, cmd.Private, cmd.CreatedAt, cmd.UpdatedAt)
}

func (f *Machine) applyReorder(raw []byte) interface{} {
	var cmd ReorderCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.ReorderThread(context.TODO(), cmd.Id, cmd.UserId, cmd.ParentId, cmd.NextId, cmd.PrevId, cmd.UpdatedAt)
}

func (f *Machine) applyUpdate(raw []byte) interface{} {
	var cmd UpdateCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.UpdateThread(context.TODO(), cmd.Id, cmd.UserId, cmd.Name, cmd.Description, cmd.UpdatedAt)
}

func (f *Machine) applyDelete(raw []byte) interface{} {
	var cmd DeleteCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.DeleteThread(context.TODO(), cmd.Id, cmd.UserId)
}

func (f *Machine) applyPublish(raw []byte) interface{} {
	var cmd PublishCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.PublishThread(context.TODO(), cmd.Id, cmd.UserId, cmd.UpdatedAt)
}

func (f *Machine) applyPrivate(raw []byte) interface{} {
	var cmd PrivateCommand
	proto.Unmarshal(raw, &cmd)

	return f.threadsRepo.PrivateThread(context.TODO(), cmd.Id, cmd.UserId, cmd.UpdatedAt)
}
