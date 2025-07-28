package users

import (
	"time"
	"errors"
	"context"

	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/users/internal/controller"
	"github.com/bd878/gallery/server/users/internal/repository"
	sessionsmodel "github.com/bd878/gallery/server/sessions/pkg/model"
)

type Repository interface {
	Add(ctx context.Context, id int32, name, password string) error
	Delete(ctx context.Context, id int32) error
	Find(ctx context.Context, params *model.FindUserParams) (*model.User, error)
	Get(ctx context.Context, id int32) (*model.User, error)
}

type SessionsGateway interface {
	GetSession(ctx context.Context, token string) (session *sessionsmodel.Session, err error)
	ListUserSessions(ctx context.Context, userID int32) (sessions []*sessionsmodel.Session, err error)
	RemoveAllUserSessions(ctx context.Context, userID int32) (err error)
	CreateSession(ctx context.Context, userID int32) (session *sessionsmodel.Session, err error)
	RemoveSession(ctx context.Context, token string) (err error)
}

type MessagesGateway interface {
	DeleteAllUserMessages(ctx context.Context, userID int32) error
}

type Controller struct {
	repo       Repository
	sessions   SessionsGateway
	messages   MessagesGateway
}

func New(repo Repository, messages MessagesGateway, sessions SessionsGateway) *Controller {
	return &Controller{repo: repo, messages: messages, sessions: sessions}
}

func (c *Controller) CreateUser(ctx context.Context, name, password string) (user *model.User, err error) {
	id := utils.RandomID()

	var session *sessionsmodel.Session
	session, err = c.sessions.CreateSession(ctx, id)
	if err != nil {
		return
	}

	err = c.repo.Add(ctx, id, name, password)
	if err != nil {
		return
	}

	user = &model.User{
		ID:                id,
		Name:              name,
		Password:          password,
		Token:             session.Token,
		ExpiresUTCNano:    session.ExpiresUTCNano,
	}

	return
}

func (c *Controller) DeleteUser(ctx context.Context, userID int32) (err error) {
	err = c.sessions.RemoveAllUserSessions(ctx, userID)
	if err != nil {
		return
	}

	err = c.messages.DeleteAllUserMessages(ctx, userID)
	if err != nil {
		return
	}

	err = c.repo.Delete(ctx, userID)
	if err != nil {
		if errors.Is(repository.ErrNoRows, err) {
			return controller.ErrUserNotFound
		}
	}
	return
}

func (c *Controller) LogoutUser(ctx context.Context, token string) (err error) {
	err = c.sessions.RemoveSession(ctx, token)
	return
}

func (c *Controller) FindUser(ctx context.Context, params *model.FindUserParams) (user *model.User, err error) {
	if params.Token != "" {
		session, err := c.sessions.GetSession(ctx, params.Token)
		if err != nil {
			return nil, err
		}

		user, err = c.repo.Get(ctx, session.UserID)
	} else {
		user, err = c.repo.Find(ctx, params)
	}

	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, controller.ErrUserNotFound
		} else {
			return nil, err
		}
	}

	return
}

func (c *Controller) LoginUser(ctx context.Context, name, password string) (session *sessionsmodel.Session, err error) {
	user, err := c.repo.Find(ctx, &model.FindUserParams{
		Name: name,
	})
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, controller.ErrUserNotFound
		}
		return nil, err
	}

	if user.Password != password {
		return nil, controller.ErrWrongPassword
	}

	session, err = c.sessions.CreateSession(ctx, user.ID)

	return
}

func (c *Controller) AuthUser(ctx context.Context, token string) (user *model.User, err error) {
	session, err := c.sessions.GetSession(ctx, token)
	if err != nil {
		return nil, err
	}

	if session.ExpiresUTCNano <= time.Now().UnixNano() {
		return nil, controller.ErrTokenExpired
	}

	user, err = c.repo.Get(ctx, session.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil, controller.ErrUserNotFound
		}

		return
	}

	user.Token = session.Token
	user.ExpiresUTCNano = session.ExpiresUTCNano

	return
}

func (c *Controller) GetUser(ctx context.Context, id int32) (user *model.User, err error) {
	user, err = c.repo.Get(ctx, id)
	if errors.Is(err, repository.ErrNoRows) {
		return nil, controller.ErrUserNotFound
	}
	return
}
