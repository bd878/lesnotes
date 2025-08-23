package users

import (
	"time"
	"context"
	"golang.org/x/crypto/bcrypt"

	"github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/users/internal/controller"
	sessions "github.com/bd878/gallery/server/sessions/pkg/model"
)

type Repository interface {
	Save(ctx context.Context, id int64, login, password string) error
	Delete(ctx context.Context, id int64) error
	Find(ctx context.Context, id int64, login string) (*model.User, error)
}

type SessionsGateway interface {
	GetSession(ctx context.Context, token string) (session *sessions.Session, err error)
	ListUserSessions(ctx context.Context, userID int64) (sessions []*sessions.Session, err error)
	RemoveAllUserSessions(ctx context.Context, userID int64) (err error)
	CreateSession(ctx context.Context, userID int64) (session *sessions.Session, err error)
	RemoveSession(ctx context.Context, token string) (err error)
}

type MessagesGateway interface {
	DeleteAllUserMessages(ctx context.Context, userID int64) error
}

type Controller struct {
	repo       Repository
	sessions   SessionsGateway
	messages   MessagesGateway
}

func New(repo Repository, messages MessagesGateway, sessions SessionsGateway) *Controller {
	return &Controller{repo: repo, messages: messages, sessions: sessions}
}

func (c *Controller) CreateUser(ctx context.Context, id int64, login, password string) (user *model.User, err error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	session, err := c.sessions.CreateSession(ctx, id)
	if err != nil {
		return nil, err
	}

	err = c.repo.Save(ctx, id, login, string(hashed))
	if err != nil {
		return
	}

	user = &model.User{
		ID:                id,
		Login:             login,
		HashedPassword:    string(hashed),
		Token:             session.Token,
		ExpiresUTCNano:    session.ExpiresUTCNano,
	}

	return
}

func (c *Controller) DeleteUser(ctx context.Context, userID int64) (err error) {
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
		return controller.ErrUserNotFound
	}
	return
}

func (c *Controller) LogoutUser(ctx context.Context, token string) (err error) {
	err = c.sessions.RemoveSession(ctx, token)
	return
}

func (c *Controller) FindUser(ctx context.Context, id int64, login, token string) (user *model.User, err error) {
	if token != "" {
		session, err := c.sessions.GetSession(ctx, token)
		if err != nil {
			return nil, err
		}

		user, err = c.repo.Find(ctx, int64(session.UserID), "")
	} else {
		user, err = c.repo.Find(ctx, 0, login)
	}

	if err != nil {
		return nil, err
	}

	return
}

func (c *Controller) LoginUser(ctx context.Context, login, hashedPassword string) (session *sessions.Session, err error) {
	user, err := c.repo.Find(ctx, 0, login)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(hashedPassword))
	if err != nil {
		return nil, err
	}

	session, err = c.sessions.CreateSession(ctx, int64(user.ID))

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

	user, err = c.repo.Find(ctx, int64(session.UserID), "")
	if err != nil {
		return nil, controller.ErrUserNotFound
	}

	user.Token = session.Token
	user.ExpiresUTCNano = session.ExpiresUTCNano

	return
}

func (c *Controller) GetUser(ctx context.Context, id int64) (user *model.User, err error) {
	user, err = c.repo.Find(ctx, id, "")
	if err != nil {
		return nil, controller.ErrUserNotFound
	}
	return
}
