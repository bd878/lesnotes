package machine

import (
	"io"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/internal/logger"
)

type UsersRepository interface {
	Save(ctx context.Context, id int64, login, hashedPassword string, metadata []byte) (err error)
	Delete(ctx context.Context, id int64) (err error)
	Update(ctx context.Context, id int64, newLogin string, metadata []byte) (err error)
	Find(ctx context.Context, id int64/*TODO: not used, remove*/, login string) (*model.User, error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

var _ raft.FSM = (*Machine)(nil)

type Machine struct {
	log         *logger.Logger
	usersRepo   UsersRepository
}

func New(usersRepo UsersRepository, log *logger.Logger) *Machine {
	return &Machine{
		log:         log,
		usersRepo:   usersRepo,
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
		logger.Errorw("unknown request type", "type", reqType)
	}
	return nil
}

func (f *Machine) applyAppend(raw []byte) interface{} {
	var cmd AppendCommand
	proto.Unmarshal(raw, &cmd)

	return f.usersRepo.Save(context.Background(), cmd.Id, cmd.Login, cmd.HashedPassword, cmd.Metadata)
}

func (f *Machine) applyUpdate(raw []byte) interface{} {
	var cmd UpdateCommand
	proto.Unmarshal(raw, &cmd)

	return f.usersRepo.Update(context.Background(), cmd.Id, cmd.Login, cmd.Metadata)
}

func (f *Machine) applyDelete(raw []byte) interface{} {
	var cmd DeleteCommand
	proto.Unmarshal(raw, &cmd)

	return f.usersRepo.Delete(context.Background(), cmd.Id)
}
