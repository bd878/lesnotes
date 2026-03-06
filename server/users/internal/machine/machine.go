package machine

import (
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
)

type UsersRepository interface {
	Save(ctx context.Context, id int64, login, hashedPassword string, metadata []byte, createdAt, updatedAt string) (err error)
	Delete(ctx context.Context, id int64) (err error)
	Update(ctx context.Context, id int64, login *string, metadata []byte, updatedAt string) (err error)
}

type UsersDumper interface {
	Open(ctx context.Context) (ch chan *model.User, err error)
	Restore(ctx context.Context, user *model.User) (err error)
	Close() (err error)
}

var _ raft.FSM = (*Machine)(nil)

type Machine struct {
	log         *logger.Logger
	usersRepo   UsersRepository
	usersDumper UsersDumper
}

func New(usersRepo UsersRepository, usersDumper UsersDumper, log *logger.Logger) *Machine {
	return &Machine{
		log:         log,
		usersRepo:   usersRepo,
		usersDumper: usersDumper,
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
	default:
		f.log.Errorw("unknown request type", "type", reqType)
	}
	return nil
}

func (f *Machine) applyAppend(raw []byte) interface{} {
	var cmd AppendCommand
	proto.Unmarshal(raw, &cmd)

	return f.usersRepo.Save(context.TODO(), cmd.Id, cmd.Login, cmd.HashedPassword, cmd.Metadata, cmd.CreatedAt, cmd.UpdatedAt)
}

func (f *Machine) applyUpdate(raw []byte) interface{} {
	var cmd UpdateCommand
	proto.Unmarshal(raw, &cmd)

	return f.usersRepo.Update(context.TODO(), cmd.Id, cmd.Login, cmd.Metadata, cmd.UpdatedAt)
}

func (f *Machine) applyDelete(raw []byte) interface{} {
	var cmd DeleteCommand
	proto.Unmarshal(raw, &cmd)

	return f.usersRepo.Delete(context.TODO(), cmd.Id)
}
