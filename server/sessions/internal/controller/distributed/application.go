package application

import (
	"time"
	"context"
	"bytes"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/utils"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/sessions/pkg/model"
	"github.com/bd878/gallery/server/sessions/internal/machine"
)

type SessionsRepository interface {
	Get(ctx context.Context, token string) (session *model.Session, err error)
	List(ctx context.Context, userID int64) (sessions []*model.Session, err error)
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

func (m *Distributed) CreateSession(ctx context.Context, userID int64) (session *model.Session, err error) {
	m.log.Debugw("create session", "userID", userID)

	token := utils.RandomString(10)
	expiresUTCNano := time.Now().Add(time.Hour * 24 * 5).UnixNano()

	cmd, err := proto.Marshal(&machine.AppendCommand{
		UserId: userID,
		Token:  token,
		ExpiresUtcNano: expiresUTCNano,
	})
	if err != nil {
		return nil, err
	}

	err = m.apply(ctx, machine.AppendRequest, cmd)

	session = &model.Session{
		UserID:         userID,
		Token:          token,
		ExpiresUTCNano: expiresUTCNano,
	}

	return
}

func (m *Distributed) GetSession(ctx context.Context, token string) (*model.Session, error) {
	m.log.Debugw("get session", "token", token)

	return m.sessionsRepo.Get(ctx, token)
}

func (m *Distributed) ListUserSessions(ctx context.Context, userID int64) ([]*model.Session, error) {
	m.log.Debugw("list user sessions", "userID", userID)

	return m.sessionsRepo.List(ctx, userID)
}

func (m *Distributed) RemoveSession(ctx context.Context, token string) error {
	m.log.Debugw("remove session", "token", token)

	cmd, err := proto.Marshal(&machine.DeleteCommand{
		Token: token,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.DeleteRequest, cmd)

	return err
}

func (m *Distributed) RemoveUserSessions(ctx context.Context, userID int64) error {
	m.log.Debugw("remove user sessions", "userID", userID)

	cmd, err := proto.Marshal(&machine.DeleteUserSessionsCommand{
		UserId: userID,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.DeleteUserSessionsRequest, cmd)

	return err
}

func (m *Distributed) GetServers(ctx context.Context) ([]*api.Server, error) {
	m.log.Debugln("get servers")
	return m.consensus.GetServers(ctx)
}
