package service

import (
	"fmt"
	"time"
	"context"
	"golang.org/x/crypto/bcrypt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/api/users"
	"github.com/bd878/gallery/server/internal/rpc"
	"github.com/bd878/gallery/server/internal/ddd"
	"github.com/bd878/gallery/server/users/config"
	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/db/users/pkg/loadbalance"
	"github.com/bd878/gallery/server/db/users/pkg/machine"
	"github.com/bd878/gallery/server/users/internal/domain"
	"github.com/bd878/gallery/server/users/internal/controller"
	sessions "github.com/bd878/gallery/server/db/sessions/pkg/model"
)

type SessionsGateway interface {
	GetSession(ctx context.Context, token string) (session *sessions.Session, err error)
	ListUserSessions(ctx context.Context, userID int64) (sessions []*sessions.Session, err error)
	RemoveAllUserSessions(ctx context.Context, userID int64) (err error)
	CreateSession(ctx context.Context, userID int64) (session *sessions.Session, err error)
	RemoveSession(ctx context.Context, token string) (err error)
}

type Controller struct {
	conf         config.Config
	client       users.UsersClient
	conn         *grpc.ClientConn
	sessions     SessionsGateway
	publisher    ddd.EventPublisher[ddd.Event]
}

func New(conf config.Config, sessions SessionsGateway, publisher ddd.EventPublisher[ddd.Event]) *Controller {
	controller := &Controller{
		conf: conf,
		sessions: sessions,
		publisher: publisher,
	}

	controller.setupConnection()

	return controller
}

func (s *Controller) Close() {
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			logger.Error(zap.Error(err))
		}
	}
}

func (s *Controller) setupConnection() (err error) {
	s.Close()

	conn, err := rpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			s.conf.UsersServiceAddr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	client := users.NewUsersClient(conn)

	s.conn = conn
	s.client = client

	return
}

func (s *Controller) isConnFailed() bool {
	state := s.conn.GetState()
	if state == connectivity.Shutdown ||
		state == connectivity.TransientFailure ||
		state == connectivity.Connecting {
		logger.Debugw("users conn failed", "state", state.String())
		return true
	}
	return false
}

func (s *Controller) CreateUser(ctx context.Context, id int64, login, password string) (user *model.User, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("create user", "id", id, "login", login, "len(password)", len(password))

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	cmd, err := proto.Marshal(&users.AppendCommand{
		Id:             id,
		Login:          login,
		HashedPassword: string(hashed),
		CreatedAt:      time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:      time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.AppendRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return
	}

	logger.Debugln("user created")
	logger.Debugln("create session")

	session, err := s.sessions.CreateSession(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Debugln("session created")

	user = &model.User{
		ID:        id,
		Login:     login,
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt,
	}

	return
}

func (s *Controller) FindUser(ctx context.Context, id int64, login, token string) (user *model.User, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("find user", "id", id, "login", login, "token", token)

	var userProto *users.User
	if token != "" {
		session, err := s.sessions.GetSession(ctx, token)
		if err != nil {
			return nil, err
		}

		userProto, err = s.client.GetUser(ctx, &users.GetUserRequest{
			Id: int64(session.UserID),
		})
	} else {
		userProto, err = s.client.FindUser(ctx, &users.FindUserRequest{
			Login: login,
		})
	}
	if err != nil {
		return
	}

	user = model.UserFromProto(userProto)

	return
}

func (s *Controller) AuthUser(ctx context.Context, token string) (user *model.User, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("auth user", "token", token)

	session, err := s.sessions.GetSession(ctx, token)
	if err != nil {
		return nil, err
	}

	expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	if expiresAt.Unix() <= time.Now().Unix() {
		return nil, controller.ErrTokenExpired
	}

	var userProto *users.User
	userProto, err = s.client.GetUser(ctx, &users.GetUserRequest{
		Id:  int64(session.UserID),
	})
	if err != nil {
		return
	}

	user = model.UserFromProto(userProto)

	user.Token = session.Token
	user.ExpiresAt = session.ExpiresAt

	return
}

func (s *Controller) GetUser(ctx context.Context, id int64) (user *model.User, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("get user", "id", id)

	userProto, err := s.client.GetUser(ctx, &users.GetUserRequest{Id: id})
	if err != nil {
		return nil, err
	}

	user = model.UserFromProto(userProto)

	return
}

func (s *Controller) UpdateUser(ctx context.Context, id int64, login *string, metadata []byte) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("update user", "id", id, "login", login, "metadata", metadata)

	cmd, err := proto.Marshal(&users.UpdateCommand{
		Id:             id,
		Login:          login,
		Metadata:       metadata,
		UpdatedAt:      time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.UpdateRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) MakePremium(ctx context.Context, id int64, invoiceID, createdAt, expiresAt string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("make premium", "id", id, "invoice_id", invoiceID, "created_at", createdAt, "expiresAt", expiresAt)

	cmd, err := proto.Marshal(&users.MakePremiumCommand{
		InvoiceId:       invoiceID,
		Id:              id,
		CreatedAt:       createdAt,
		ExpiresAt:       expiresAt,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.MakePremiumRequest),
		Cmd: cmd,
		Duration: "10s",
	})

	return
}

func (s *Controller) LoginUser(ctx context.Context, login, password string) (session *sessions.Session, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("login user", "login", login, "len(password)", len(password))

	user, err := s.client.FindUser(ctx, &users.FindUserRequest{
		Login:   login,
	})
	if err != nil {
		return nil, err
	}

	logger.Debugw("user found", "user_id", user.Id)

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		logger.Errorln(err)
		return nil, controller.ErrWrongPassword
	}

	session, err = s.sessions.CreateSession(ctx, int64(user.Id))

	logger.Debugw("session created", "session", session)

	return
}

func (s *Controller) DeleteUser(ctx context.Context, id int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("delete user", "id", id)

	// TODO: emit event, not call
	err = s.sessions.RemoveAllUserSessions(ctx, id)
	if err != nil {
		return
	}

	cmd, err := proto.Marshal(&users.DeleteCommand{
		Id: id,
	})
	if err != nil {
		return err
	}

	_, err = s.client.Apply(ctx, &api.Command{
		ReqType: int32(machine.DeleteRequest),
		Cmd: cmd,
		Duration: "10s",
	})
	if err != nil {
		return err
	}

	event, err := domain.DeleteUser(id)
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, event)
}

func (s *Controller) LogoutUser(ctx context.Context, token string) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("logout user", "token", token)

	err = s.sessions.RemoveSession(ctx, token)

	return
}

