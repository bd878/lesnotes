package distributed

import (
	"io"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	usersmodel "github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/internal/logger"
)

type RepoConnection interface {
	Release()
}

type Repository interface {
	Save(ctx context.Context, id int64, login, hashedPassword string, metadata []byte) (err error)
	Delete(ctx context.Context, id int64) (err error)
	Update(ctx context.Context, id int64, newLogin string, metadata []byte) (err error)
	Find(ctx context.Context, id int64/*TODO: not used, remove*/, login string) (*usersmodel.User, error)
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
	default:
		logger.Errorw("unknown request type", "type", reqType)
	}
	return nil
}

func (f *fsm) applyAppend(raw []byte) interface{} {
	var cmd AppendCommand
	proto.Unmarshal(raw, &cmd)

	return f.repo.Save(context.Background(), cmd.Id, cmd.Login, cmd.HashedPassword, cmd.Metadata)
}

func (f *fsm) applyUpdate(raw []byte) interface{} {
	var cmd UpdateCommand
	proto.Unmarshal(raw, &cmd)

	return f.repo.Update(context.Background(), cmd.Id, cmd.Login, cmd.Metadata)
}

func (f *fsm) applyDelete(raw []byte) interface{} {
	var cmd DeleteCommand
	proto.Unmarshal(raw, &cmd)

	return f.repo.Delete(context.Background(), cmd.Id)
}
