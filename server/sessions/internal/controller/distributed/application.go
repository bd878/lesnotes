package application

import (
	"time"
	"context"
	"bytes"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/sessions"
	"github.com/bd878/gallery/server/internal/utils"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/sessions/pkg/machine"
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
	log               *logger.Logger
	sessionsRepo      SessionsRepository
}

func New(consensus Consensus, sessionsRepo SessionsRepository, log *logger.Logger) *Distributed {
	return &Distributed{
		log:            log,
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
	m.log.Debugw("create session", "userID", userID)

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
	m.log.Debugw("get session", "token", token)

	return m.sessionsRepo.Get(ctx, token)
}

func (m *Distributed) ListUserSessions(ctx context.Context, userID int64) ([]*sessions.Session, error) {
	m.log.Debugw("list user sessions", "userID", userID)

	return m.sessionsRepo.List(ctx, userID)
}

func (m *Distributed) RemoveSession(ctx context.Context, token string) error {
	m.log.Debugw("remove session", "token", token)

	cmd, err := proto.Marshal(&sessions.DeleteCommand{
		Token: token,
	})
	if err != nil {
		return err
	}

	return m.apply(ctx, machine.DeleteRequest, cmd)
}

func (m *Distributed) RemoveUserSessions(ctx context.Context, userID int64) error {
	m.log.Debugw("remove user sessions", "userID", userID)

	cmd, err := proto.Marshal(&sessions.DeleteUserSessionsCommand{
		UserId: userID,
	})
	if err != nil {
		return err
	}

	return m.apply(ctx, machine.DeleteUserSessionsRequest, cmd)
}

func (m *Distributed) GetServers(ctx context.Context) ([]*api.Server, error) {
	m.log.Debugln("get servers")
	return m.consensus.GetServers(ctx)
}
