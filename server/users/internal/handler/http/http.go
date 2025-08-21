package http

import (
	"net/http"
	"time"
	"io"
	"context"
	"encoding/json"

	"github.com/bd878/gallery/server/users/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
	sessionsmodel "github.com/bd878/gallery/server/sessions/pkg/model"
	"github.com/bd878/gallery/server/utils"
)

type Controller interface {
	CreateUser(ctx context.Context, id int64, name, password string) (user *model.User, err error)
	FindUser(ctx context.Context, params *model.FindUserParams) (user *model.User, err error)
	AuthUser(ctx context.Context, token string) (user *model.User, err error)
	GetUser(ctx context.Context, userID int32) (user *model.User, err error)
	LoginUser(ctx context.Context, name, password string) (session *sessionsmodel.Session, err error)
	DeleteUser(ctx context.Context, userID int32) (err error)
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

func getTextField(w http.ResponseWriter, req *http.Request, field string) (value string, ok bool) {
	value = req.PostFormValue(field)
	if value == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no " + field,
		})
	} else {
		ok = true
	}
	return
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
