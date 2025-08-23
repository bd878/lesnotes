package http

import (
	"net/http"
	"time"
	"io"
	"context"

	"github.com/bd878/gallery/server/utils"
	users "github.com/bd878/gallery/server/users/pkg/model"
	sessions "github.com/bd878/gallery/server/sessions/pkg/model"
)

type Controller interface {
	CreateUser(ctx context.Context, id int64, login, password string) (user *users.User, err error)
	FindUser(ctx context.Context, id int64, login, token string) (user *users.User, err error)
	AuthUser(ctx context.Context, token string) (user *users.User, err error)
	GetUser(ctx context.Context, userID int64) (user *users.User, err error)
	LoginUser(ctx context.Context, login, hashedPassword string) (session *sessions.Session, err error)
	DeleteUser(ctx context.Context, userID int64) (err error)
	LogoutUser(ctx context.Context, token string) (err error)
}

type Config struct {
	CookieDomain string
}

type Handler struct {
	controller         Controller
	config             Config
}

func New(controller Controller, config Config) *Handler {
	return &Handler{controller, config}
}

func (h *Handler) Status(w http.ResponseWriter, _ *http.Request) error {
	if _, err := io.WriteString(w, "ok\n"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}

func createToken(w http.ResponseWriter, cookieDomain string) (token string, expires int64) {
	token = utils.RandomString(10)
	expiresAt := time.Now().Add(time.Hour * 24 * 5)

	http.SetCookie(w, &http.Cookie{
		Name:             "token",
		Value:             token,
		Domain:            cookieDomain,
		Expires:           expiresAt,
		Path:             "/",
		HttpOnly:          true,
	})

	return token, expiresAt.UnixNano()
}
