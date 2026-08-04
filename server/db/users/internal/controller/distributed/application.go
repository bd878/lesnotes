package application

import (
	"time"
	"context"
	"bytes"
	"log/slog"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/users"
	"github.com/bd878/gallery/server/db/users/pkg/machine"
)

type UsersRepository interface {
	FindByID(ctx context.Context, id int64) (user *users.User, err error)
	FindByLogin(ctx context.Context, login string) (user *users.User, err error)
}

type Consensus interface {
	Apply(cmd []byte, timeout time.Duration) (err error)
	GetServers(ctx context.Context) ([]*api.Server, error)
}

type Distributed struct {
	consensus         Consensus
	usersRepo         UsersRepository
}

func New(consensus Consensus, usersRepo UsersRepository) *Distributed {
	return &Distributed{
		consensus:   consensus,
		usersRepo:   usersRepo,
	}
}

func (m *Distributed) Apply(ctx context.Context, reqType machine.RequestType, cmd []byte, duration time.Duration) (err error) {
	var buf bytes.Buffer
	_, err = buf.Write([]byte{byte(reqType)})
	if err != nil {
		return
	}

	_, err = buf.Write(cmd)
	if err != nil {
		return
	}

	return m.consensus.Apply(buf.Bytes(), duration)
}

func (m *Distributed) FindUser(ctx context.Context, login string) (user *users.User, err error) {
	slog.Debug("find user", slog.String("login", login))
	return m.usersRepo.FindByLogin(ctx, login)
}

func (m *Distributed) GetUser(ctx context.Context, id int64) (user *users.User, err error) {
	slog.Debug("get user", slog.Int64("id", id))
	return m.usersRepo.FindByID(ctx, id)
}

func (m *Distributed) GetServers(ctx context.Context) ([]*api.Server, error) {
	slog.Debug("get servers")
	return m.consensus.GetServers(ctx)
}
