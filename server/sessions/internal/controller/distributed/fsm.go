package distributed

import (
	"io"
	"context"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/sessions/pkg/model"
	"github.com/bd878/gallery/server/logger"
)

type RepoConnection interface {
	Release()
}

type Repository interface {
	Save(ctx context.Context, userID int64, token string, expiresUTCNano int64) (err error)
	Get(ctx context.Context, token string) (session *model.Session, err error)
	List(ctx context.Context, userID int64) (sessions []*model.Session, err error)
	Delete(ctx context.Context, token string) (err error)
	DeleteAll(ctx context.Context, userID int64) (err error)

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
	case DeleteRequest:
		return f.applyDelete(buf[1:])
	case DeleteUserSessionsRequest:
		return f.applyDeleteUserSessions(buf[1:])
	default:
		logger.Errorw("unknown request type", "type", reqType)
	}
	return nil
}

func (f *fsm) applyAppend(raw []byte) interface{} {
	var cmd AppendCommand
	proto.Unmarshal(raw, &cmd)

	return f.repo.Save(context.Background(), cmd.UserId, cmd.Token, cmd.ExpiresUtcNano)
}

func (f *fsm) applyDelete(raw []byte) interface{} {
	var cmd DeleteCommand
	proto.Unmarshal(raw, &cmd)

	return f.repo.Delete(context.Background(), cmd.Token)
}

func (f *fsm) applyDeleteUserSessions(raw []byte) interface{} {
	var cmd DeleteUserSessionsCommand
	proto.Unmarshal(raw, &cmd)

	return f.repo.DeleteAll(context.Background(), cmd.UserId)
}