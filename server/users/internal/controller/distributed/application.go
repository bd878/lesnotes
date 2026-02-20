package application

import (
	"time"
	"context"
	"bytes"
	"golang.org/x/crypto/bcrypt"

	"google.golang.org/protobuf/proto"
	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/users/internal/machine"
)

type UsersRepository interface {
	Find(ctx context.Context, id int64, login string) (user *model.User, err error)
}

type Consensus interface {
	Apply(cmd []byte, timeout time.Duration) (err error)
	GetServers(ctx context.Context) ([]*api.Server, error)
}

type Distributed struct {
	consensus         Consensus
	log               *logger.Logger
	usersRepo         UsersRepository
}

func New(consensus Consensus, usersRepo UsersRepository, log *logger.Logger) *Distributed {
	return &Distributed{
		log:         log,
		consensus:   consensus,
		usersRepo:   usersRepo,
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

func (m *Distributed) CreateUser(ctx context.Context, id int64, login, password string, metadata []byte) (user *model.User, err error) {
	m.log.Debugw("create user", "id", id, "login", login, "password", password, "metadata", metadata)

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	cmd, err := proto.Marshal(&machine.AppendCommand{
		Id:             id,
		Login:          login,
		HashedPassword: string(hashed),
		Metadata:       metadata,
	})
	if err != nil {
		return nil, err
	}

	err = m.apply(ctx, machine.AppendRequest, cmd)

	user = &model.User{
		ID:              id,
		Login:           login,
		HashedPassword:  string(hashed),
		Metadata:        metadata,
	}

	return
}

func (m *Distributed) DeleteUser(ctx context.Context, id int64) (err error) {
	m.log.Debugw("delete user", "id", id)

	cmd, err := proto.Marshal(&machine.DeleteCommand{
		Id: id,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.DeleteRequest, cmd)
	if err != nil {
		return
	}

	return nil
}

func (m *Distributed) UpdateUser(ctx context.Context, id int64, login string, metadata []byte) (err error) {
	m.log.Debugw("update user", "id", id, "login", login, "metadata", metadata)

	cmd, err := proto.Marshal(&machine.UpdateCommand{
		Id:          id,
		Login:       login,
		Metadata:    metadata,
	})
	if err != nil {
		return err
	}

	err = m.apply(ctx, machine.UpdateRequest, cmd)

	return
}

func (m *Distributed) FindUser(ctx context.Context, login string) (user *model.User, err error) {
	m.log.Debugw("find user", "login", login)
	return m.usersRepo.Find(ctx, 0, login)
}

func (m *Distributed) GetUser(ctx context.Context, id int64) (user *model.User, err error) {
	m.log.Debugw("get user", "id", id)
	return m.usersRepo.Find(ctx, id, "")
}

func (m *Distributed) GetServers(ctx context.Context) ([]*api.Server, error) {
	m.log.Debugln("get servers")
	return m.consensus.GetServers(ctx)
}
