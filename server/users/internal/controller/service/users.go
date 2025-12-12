package service

import (
	"fmt"
	"time"
	"context"
	"golang.org/x/crypto/bcrypt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/users/pkg/loadbalance"
	"github.com/bd878/gallery/server/users/internal/controller"
	sessions "github.com/bd878/gallery/server/sessions/pkg/model"
)

type Config struct {
	RpcAddr string
}

type SessionsGateway interface {
	GetSession(ctx context.Context, token string) (session *sessions.Session, err error)
	ListUserSessions(ctx context.Context, userID int64) (sessions []*sessions.Session, err error)
	RemoveAllUserSessions(ctx context.Context, userID int64) (err error)
	CreateSession(ctx context.Context, userID int64) (session *sessions.Session, err error)
	RemoveSession(ctx context.Context, token string) (err error)
}

type MessagesGateway interface {
	DeleteUserMessages(ctx context.Context, userID int64) error
}

type Controller struct {
	conf         Config
	client       api.UsersClient
	conn         *grpc.ClientConn
	sessions     SessionsGateway
	messages     MessagesGateway
}

func New(conf Config, messages MessagesGateway, sessions SessionsGateway) *Controller {
	controller := &Controller{conf: conf, messages: messages, sessions: sessions}

	controller.setupConnection()

	return controller
}

func (s *Controller) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *Controller) setupConnection() (err error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:///%s",
			loadbalance.Name,
			s.conf.RpcAddr,
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	client := api.NewUsersClient(conn)

	s.conn = conn
	s.client = client

	return
}

func (s *Controller) isConnFailed() bool {
	state := s.conn.GetState()
	if state == connectivity.Shutdown || state == connectivity.TransientFailure {
		logger.Debugln("connection failed")
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

	_, err = s.client.CreateUser(ctx, &api.CreateUserRequest{
		Id:         id,
		Login:      login,
		Password:   password,
	})
	if err != nil {
		return
	}

	session, err := s.sessions.CreateSession(ctx, id)
	if err != nil {
		return nil, err
	}

	user = &model.User{
		ID:       id,
		Login:    login,
		Token:    session.Token,
		ExpiresUTCNano: session.ExpiresUTCNano,
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

	var userProto *api.User
	if token != "" {
		session, err := s.sessions.GetSession(ctx, token)
		if err != nil {
			return nil, err
		}

		userProto, err = s.client.GetUser(ctx, &api.GetUserRequest{
			Id: int64(session.UserID),
		})
	} else {
		userProto, err = s.client.FindUser(ctx, &api.FindUserRequest{
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

	if session.ExpiresUTCNano <= time.Now().UnixNano() {
		return nil, controller.ErrTokenExpired
	}

	var userProto *api.User
	userProto, err = s.client.GetUser(ctx, &api.GetUserRequest{
		Id:  int64(session.UserID),
	})
	if err != nil {
		return
	}

	user = model.UserFromProto(userProto)

	user.Token = session.Token
	user.ExpiresUTCNano = session.ExpiresUTCNano

	return
}

func (s *Controller) GetUser(ctx context.Context, id int64) (user *model.User, err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("get user", "id", id)

	userProto, err := s.client.GetUser(ctx, &api.GetUserRequest{Id: id})
	if err != nil {
		return nil, err
	}

	user = model.UserFromProto(userProto)

	return
}

func (s *Controller) UpdateUser(ctx context.Context, id int64, newLogin string, metadata []byte) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("update user", "id", id, "login", newLogin, "metadata", metadata)

	_, err = s.client.UpdateUser(ctx, &api.UpdateUserRequest{
		Id:        id,
		Login:     newLogin,
		Metadata:  metadata,
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

	user, err := s.client.FindUser(ctx, &api.FindUserRequest{
		Login:   login,
	})
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		logger.Errorln(err)
		return nil, controller.ErrWrongPassword
	}

	session, err = s.sessions.CreateSession(ctx, int64(user.Id))

	return
}

func (s *Controller) DeleteUser(ctx context.Context, id int64) (err error) {
	if s.isConnFailed() {
		if err = s.setupConnection(); err != nil {
			return
		}
	}

	logger.Debugw("delete user", "id", id)

	err = s.sessions.RemoveAllUserSessions(ctx, id)
	if err != nil {
		return
	}

	err = s.messages.DeleteUserMessages(ctx, id)
	if err != nil {
		return
	}

	_, err = s.client.DeleteUser(ctx, &api.DeleteUserRequest{
		Id: id,
	})

	return
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

