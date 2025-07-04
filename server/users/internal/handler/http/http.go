package http

import (
	"net/http"
	"time"
	"io"
	"context"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
	messagesmodel "github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/utils"
)

type Controller interface {
	AddUser(ctx context.Context, log *logger.Logger, params *model.AddUserParams) (int32, error)
	HasUser(ctx context.Context, log *logger.Logger, params *model.HasUserParams) (bool, error)
	RefreshToken(ctx context.Context, log *logger.Logger, params *model.RefreshTokenParams) error
	DeleteToken(ctx context.Context, log *logger.Logger, params *model.DeleteTokenParams) error
	GetUser(ctx context.Context, log *logger.Logger, params *model.GetUserParams) (*model.User, error)
	DeleteUser(ctx context.Context, log *logger.Logger, params *model.DeleteUserParams) error
}

type MessagesGateway interface {
	DeleteAllUserMessages(ctx context.Context, log *logger.Logger, params *messagesmodel.DeleteAllUserMessagesParams) error
}

type Config struct {
	CookieDomain string
}

type Handler struct {
	controller      Controller
	config          Config
	messagesGateway MessagesGateway
}

func New(controller Controller, config Config, messagesGateway MessagesGateway) *Handler {
	return &Handler{controller, config, messagesGateway}
}

func (h *Handler) Status(log *logger.Logger, w http.ResponseWriter, _ *http.Request) error {
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

func refreshToken(h *Handler, w http.ResponseWriter, req *http.Request, userName string) (*model.User, error) {
	token, expiresUtcNano := createToken(w, h.config.CookieDomain)

	user := &model.User{
		Name:               userName,
		Token:              token,
		ExpiresUTCNano:     expiresUtcNano,
	}

	err := h.controller.RefreshToken(req.Context(), logger.Default(), &model.RefreshTokenParams{
		User: user,
	})
	return user, err
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

func attachTokenToResponse(w http.ResponseWriter, token, cookieDomain string, expiresUtcNano int64) {
	http.SetCookie(w, &http.Cookie{
		Name:          "token",
		Value:          token,
		Domain:         cookieDomain,
		Expires:        time.Unix(0, expiresUtcNano),
		Path:          "/",
		HttpOnly:       true,
	})
}

func deleteToken(w http.ResponseWriter, cookieDomain string) {
	http.SetCookie(w, &http.Cookie{
		Name:           "token",
		Value:          "",
		Domain:         cookieDomain,
		Expires:        time.Unix(0, 0),
		Path: "/",
		HttpOnly: true,
	})
}