package users

import (
	"time"
	"context"
	"golang.org/x/crypto/bcrypt"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/users/internal/controller"
	sessions "github.com/bd878/gallery/server/sessions/pkg/model"
)

type Repository interface {
	Save(ctx context.Context, id int64, login, salt, theme, lang string, fontSize int32) (err error)
	Delete(ctx context.Context, id int64) (err error)
	Find(ctx context.Context, id int64, login string) (*model.User, error)
	Update(ctx context.Context, id int64, newLogin, newTheme, newLang string, newFontSize int32) (err error)
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
	repo       Repository
	sessions   SessionsGateway
	messages   MessagesGateway
}

func New(repo Repository, messages MessagesGateway, sessions SessionsGateway) *Controller {
	return &Controller{repo: repo, messages: messages, sessions: sessions}
}

func (c *Controller) CreateUser(ctx context.Context, id int64, login, password string) (user *model.User, err error) {
	logger.Debugw("create user", "id", id, "login", login, "password", password)

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	session, err := c.sessions.CreateSession(ctx, id)
	if err != nil {
		return nil, err
	}

	err = c.repo.Save(ctx, id, login, string(hashed), "light", "", 0)
	if err != nil {
		return
	}

	user = &model.User{
		ID:                id,
		Login:             login,
		Theme:             "light",
		Lang:              "",
		HashedPassword:    string(hashed),
		Token:             session.Token,
		ExpiresUTCNano:    session.ExpiresUTCNano,
		FontSize:          0,
	}

	return
}

func (c *Controller) DeleteUser(ctx context.Context, userID int64) (err error) {
	logger.Debugw("delete user", "id", userID)

	err = c.sessions.RemoveAllUserSessions(ctx, userID)
	if err != nil {
		return
	}

	err = c.messages.DeleteUserMessages(ctx, userID)
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
	logger.Debugw("logout user", "token", token)

	err = c.sessions.RemoveSession(ctx, token)

	return
}

func (c *Controller) FindUser(ctx context.Context, id int64, login, token string) (user *model.User, err error) {
	logger.Debugw("find user", "id", id, "login", login, "token", token)

	if token != "" {
		session, err := c.sessions.GetSession(ctx, token)
		if err != nil {
			return nil, err
		}

		user, err = c.repo.Find(ctx, int64(session.UserID), "")
	} else {
		user, err = c.repo.Find(ctx, 0, login)
	}

	return
}

func (c *Controller) LoginUser(ctx context.Context, login, password string) (session *sessions.Session, err error) {
	logger.Debugw("login user", "login", login, "password", password)

	user, err := c.repo.Find(ctx, 0, login)
	if err != nil {
		logger.Errorln(err)
		return nil, controller.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		logger.Errorln(err)
		return nil, controller.ErrWrongPassword
	}

	session, err = c.sessions.CreateSession(ctx, int64(user.ID))

	return
}

func (c *Controller) AuthUser(ctx context.Context, token string) (user *model.User, err error) {
	logger.Debugw("auth user", "token", token)

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
	logger.Debugw("get user", "id", id)

	user, err = c.repo.Find(ctx, id, "")
	if err != nil {
		return nil, controller.ErrUserNotFound
	}

	return
}

func (c *Controller) UpdateUser(ctx context.Context, id int64, newLogin, newTheme, newLang string, newFontSize int32) (err error) {
	logger.Debugw("update user", "login", newLogin, "theme", newTheme, "lang", newLang, "font_size", newFontSize)

	err = c.repo.Update(ctx, id, newLogin, newTheme, newLang, newFontSize)

	return
}