package application

import (
	"time"
	"context"
	"bytes"
	"log/slog"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/sessions"
	"github.com/bd878/gallery/server/internal/utils"
	"github.com/bd878/gallery/server/db/sessions/pkg/machine"
)

type SessionsRepository interface {
	Get(ctx context.Context, token string) (session *sessions.Session, err error)
	List(ctx context.Context, userID int64) (sessions []*sessions.Session, err error)
}

type Consensus interface {
	Apply(cmd []byte, timeout time.Duration) (err error)
	GetServers(ctx context.Context) ([]*api.Server, error)
}

type Distributed struct {
	consensus         Consensus
	sessionsRepo      SessionsRepository
}

func New(consensus Consensus, sessionsRepo SessionsRepository) *Distributed {
	return &Distributed{
		consensus:      consensus,
		sessionsRepo:   sessionsRepo,
	}
}

func (m *Distributed) apply(ctx context.Context, reqType machine.RequestType, cmd []byte) (err error) {
	var buf bytes.Buffer
	_, err = buf.Write([]byte{byte(reqType)})
	if err != nil {
		return
	}

	_, err = buf.Write(cmd)
	if err != nil {
		return
	}

	return m.consensus.Apply(buf.Bytes(), 10*time.Second)
}

func (m *Distributed) CreateSession(ctx context.Context, userID int64) (session *sessions.Session, err error) {
	slog.Debug("create session", slog.Int64("userID", userID))

	token := utils.RandomString(10)

	createdAt := time.Now().UTC().Format(time.RFC3339)
	expiresAt := time.Now().Add(time.Hour * 24 * 5).UTC().Format(time.RFC3339)

	cmd, err := proto.Marshal(&sessions.AppendCommand{
		UserId:     userID,
		Token:      token,
		ExpiresAt:  expiresAt,
		CreatedAt:  createdAt,
	})
	if err != nil {
		return nil, err
	}

	err = m.apply(ctx, machine.AppendRequest, cmd)

	session = &sessions.Session{
		UserId:         userID,
		Token:          token,
		CreatedAt:      createdAt,
		ExpiresAt:      expiresAt,
	}

	return
}

func (m *Distributed) GetSession(ctx context.Context, token string) (*sessions.Session, error) {
	slog.Debug("get session", slog.String("token", token))

	return m.sessionsRepo.Get(ctx, token)
}

func (m *Distributed) ListUserSessions(ctx context.Context, userID int64) ([]*sessions.Session, error) {
	slog.Debug("list user sessions", slog.Int64("userID", userID))

	return m.sessionsRepo.List(ctx, userID)
}

func (m *Distributed) RemoveSession(ctx context.Context, token string) error {
	slog.Debug("remove session", slog.String("token", token))

	cmd, err := proto.Marshal(&sessions.DeleteCommand{
		Token: token,
	})
	if err != nil {
		return err
	}

	return m.apply(ctx, machine.DeleteRequest, cmd)
}

func (m *Distributed) RemoveUserSessions(ctx context.Context, userID int64) error {
	slog.Debug("remove user sessions", slog.Int64("userID", userID))

	cmd, err := proto.Marshal(&sessions.DeleteUserSessionsCommand{
		UserId: userID,
	})
	if err != nil {
		return err
	}

	return m.apply(ctx, machine.DeleteUserSessionsRequest, cmd)
}

func (m *Distributed) GetServers(ctx context.Context) ([]*api.Server, error) {
	slog.Debug("get servers")
	return m.consensus.GetServers(ctx)
}
