package machine

import (
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
)

type SessionsRepository interface {
	Save(ctx context.Context, userID int64, token, createdAt, expiresAt string) (err error)
	Delete(ctx context.Context, token string) (err error)
	DeleteAll(ctx context.Context, userID int64) (err error)
}

type Dumper interface {
	Open(ctx context.Context) (ch chan *api.SessionsSnapshot, err error)
	Restore(ctx context.Context, user *api.SessionsSnapshot) (err error)
	Close() (err error)
}

var _ raft.FSM = (*Machine)(nil)

type Machine struct {
	log            *logger.Logger
	dumper         Dumper
	sessionsRepo   SessionsRepository
}

func New(sessionsRepo SessionsRepository, dumper Dumper, log *logger.Logger) *Machine {
	return &Machine{
		log:          log,
		dumper:       dumper,
		sessionsRepo: sessionsRepo,
	}
}

func (f *Machine) Apply(record *raft.Log) interface{} {
	buf := record.Data
	reqType := RequestType(buf[0])
	switch reqType {
	case AppendRequest:
		return f.applyAppend(buf[1:])
	case DeleteRequest:
		return f.applyDelete(buf[1:])
	case DeleteUserSessionsRequest:
		return f.applyDeleteUserSessions(buf[1:])
	default:
		f.log.Errorw("unknown request type", "type", reqType)
	}
	return nil
}

func (f *Machine) applyAppend(raw []byte) interface{} {
	var cmd AppendCommand
	proto.Unmarshal(raw, &cmd)

	return f.sessionsRepo.Save(context.TODO(), cmd.UserId, cmd.Token, cmd.CreatedAt, cmd.ExpiresAt)
}

func (f *Machine) applyDelete(raw []byte) interface{} {
	var cmd DeleteCommand
	proto.Unmarshal(raw, &cmd)

	return f.sessionsRepo.Delete(context.TODO(), cmd.Token)
}

func (f *Machine) applyDeleteUserSessions(raw []byte) interface{} {
	var cmd DeleteUserSessionsCommand
	proto.Unmarshal(raw, &cmd)

	return f.sessionsRepo.DeleteAll(context.TODO(), cmd.UserId)
}