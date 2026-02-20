package machine

import (
	"io"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/internal/logger"
)

type SessionsRepository interface {
	Save(ctx context.Context, userID int64, token string, expiresUTCNano int64) (err error)
	Delete(ctx context.Context, token string) (err error)
	DeleteAll(ctx context.Context, userID int64) (err error)
	Dump(ctx context.Context, writer io.Writer) (err error)
	Restore(ctx context.Context, reader io.Reader) (err error)
}

var _ raft.FSM = (*Machine)(nil)

type Machine struct {
	log            *logger.Logger
	sessionsRepo   SessionsRepository
}

func New(sessionsRepo SessionsRepository, log *logger.Logger) *Machine {
	return &Machine{
		log:          log,
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

	return f.sessionsRepo.Save(context.Background(), cmd.UserId, cmd.Token, cmd.ExpiresUtcNano)
}

func (f *Machine) applyDelete(raw []byte) interface{} {
	var cmd DeleteCommand
	proto.Unmarshal(raw, &cmd)

	return f.sessionsRepo.Delete(context.Background(), cmd.Token)
}

func (f *Machine) applyDeleteUserSessions(raw []byte) interface{} {
	var cmd DeleteUserSessionsCommand
	proto.Unmarshal(raw, &cmd)

	return f.sessionsRepo.DeleteAll(context.Background(), cmd.UserId)
}