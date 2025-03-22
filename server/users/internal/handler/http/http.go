package http

import (
	"net/http"
	"time"
	"io"
	"context"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/internal/controller"
	"github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/utils"
)

type Controller interface {
	AddUser(ctx context.Context, log *logger.Logger, params *model.AddUserParams) error
	HasUser(ctx context.Context, log *logger.Logger, params *model.HasUserParams) (bool, error)
	RefreshToken(ctx context.Context, log *logger.Logger, params *model.RefreshTokenParams) error
	DeleteToken(ctx context.Context, log *logger.Logger, params *model.DeleteTokenParams) error
	GetUser(ctx context.Context, log *logger.Logger, params *model.GetUserParams) (*model.User, error)
}

type Config struct {
	CookieDomain string
}

type Handler struct {
	controller      Controller
	config          Config
}

func New(controller Controller, config Config) *Handler {
	return &Handler{controller, config}
}

func (h *Handler) Logout(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("token")
	if err != nil {
		log.Error("bad cookie")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token := cookie.Value

	deleteToken(w, h.config.CookieDomain)

	err = h.controller.DeleteToken(context.Background(), log, &model.DeleteTokenParams{
		Token: token,
	})
	if err != nil {
		log.Errorw("failed to delete token, continue", "error", err)
	}

	json.NewEncoder(w).Encode(model.ServerResponse{
		Status: "ok",
		Description: "deleted",
	})
}

func (h *Handler) Login(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
	userName, ok := getTextField(w, req, "name")
	if !ok {
		log.Error("cannot get name")
		return
	}

	password, ok := getTextField(w, req, "password")
	if !ok {
		log.Error("cannot get password")
		return
	}

	exists, err := h.controller.HasUser(context.Background(), log, &model.HasUserParams{
		User: &model.User{
			Name: userName,
			Password: password,
		},
	})
	if err != nil {
		log.Error("failed to check if user exists:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "ok",
			Description: "no user,password pair",
		})
		return
	}

	user, err := h.controller.GetUser(context.Background(), log, &model.GetUserParams{Name: userName})
	switch err {
	case controller.ErrTokenExpired:
		log.Infoln("token expired")
		err := refreshToken(h, w, userName)
		if err != nil {
			log.Error("Cannot refresh token: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case controller.ErrNotFound:
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "ok",
			Description: "no user,password pair",
		})

	case nil:
		attachTokenToResponse(w, user.Token, h.config.CookieDomain, user.ExpiresUTCNano)

	default:
		log.Error("Unknown error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.ServerResponse{
		Status: "ok",
		Description: "authenticated",
	})
}

// TODO: Invalidate stale sessions.
// User logs in in one device, get token,
// then logs in in another device, receives new token.
// Old token invalidates, though not expired...
// Check if stage.lesnotes.space tokens influences on
// lesnotes.space (it has .lesnotes.space domain)
func (h *Handler) Auth(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("token")
	if err != nil {
		log.Error("bad cookie")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token := cookie.Value

	user, err := h.controller.GetUser(context.Background(), log, &model.GetUserParams{Token: token})
	if err == controller.ErrTokenExpired {
		log.Infoln("token expired")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ServerAuthorizeResponse{
			ServerResponse: model.ServerResponse{
				Status: "ok",
				Description: "token expired",
			},
			Expired: true,
		})
	}

	if err != nil {
		log.Errorw("failed to get user by token", "token", token, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(model.ServerAuthorizeResponse{
		ServerResponse: model.ServerResponse{
			Status: "ok",
			Description: "token valid",
		},
		Expired: false,
		User: model.User{
			ID:               user.ID,
			Name:             user.Name,
			Token:            user.Token,
			ExpiresUTCNano:   user.ExpiresUTCNano,
		},
	})
}

func (h *Handler) Signup(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
	userName, ok := getTextField(w, req, "name")
	if !ok {
		log.Error("cannot get user name")
		return
	}

	exists, err := h.controller.HasUser(context.Background(), log, &model.HasUserParams{
		User: &model.User{Name: userName},
	})
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "ok",
			Description: "name exists",
		})
		return
	}

	password, ok := getTextField(w, req, "password")
	if !ok {
		log.Error("cannot get password from request")
		return
	}

	token, expiresUtcNano := createToken(w, h.config.CookieDomain)

	log.Infoln("user, password, token, expires", userName, password, token, expiresUtcNano)
	err = h.controller.AddUser(context.Background(), log, &model.AddUserParams{
		User: &model.User{
			Name:                  userName,
			Password:              password,
			Token:                 token,
			ExpiresUTCNano:        expiresUtcNano,
		},
	})
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(model.ServerResponse{
		Status: "ok",
		Description: "created",
	})
}

func (h *Handler) Status(log *logger.Logger, w http.ResponseWriter, _ *http.Request) {
	if _, err := io.WriteString(w, "ok\n"); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getTextField(w http.ResponseWriter, req *http.Request, field string) (value string, ok bool) {
	value = req.PostFormValue(field)
	if value == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ServerResponse{
			Status: "ok",
			Description: "no " + field,
		})
	} else {
		ok = true
	}
	return
}

func refreshToken(h *Handler, w http.ResponseWriter, userName string) error {
	token, expiresUtcNano := createToken(w, h.config.CookieDomain)

	return h.controller.RefreshToken(context.Background(), logger.Default(), &model.RefreshTokenParams{
		User: &model.User{
			Name:               userName,
			Token:              token,
			ExpiresUTCNano:     expiresUtcNano,
		},
	})
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