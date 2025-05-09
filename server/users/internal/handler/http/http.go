package http

import (
	"net/http"
	"time"
	"io"
	"context"
	"errors"
	"strconv"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/internal/controller"
	"github.com/bd878/gallery/server/users/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
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

func (h *Handler) Logout(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	cookie, err := req.Cookie("token")
	if err != nil {
		log.Error("bad cookie")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "bad cookie",
		})

		return err
	}

	token := cookie.Value

	deleteToken(w, h.config.CookieDomain)

	err = h.controller.DeleteToken(req.Context(), log, &model.DeleteTokenParams{
		Token: token,
	})
	if err != nil {
		log.Errorw("failed to delete token, continue", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to delete token",
		})

		return err
	}

	json.NewEncoder(w).Encode(servermodel.ServerResponse{
		Status: "ok",
		Description: "logged out",
	})
	return nil
}

func (h *Handler) Login(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	userName, ok := getTextField(w, req, "name")
	if !ok {
		log.Error("cannot get name")
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no name",
		})

		return errors.New("no name field")
	}

	password, ok := getTextField(w, req, "password")
	if !ok {
		log.Error("cannot get password")
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no password",
		})

		return errors.New("no password field")
	}

	exists, err := h.controller.HasUser(req.Context(), log, &model.HasUserParams{
		User: &model.User{
			Name: userName,
			Password: password,
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot find user",
		})

		return err
	}

	if !exists {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no user,password pair",
		})

		return errors.New("no user,password pair")
	}

	user, err := h.controller.GetUser(req.Context(), log, &model.GetUserParams{Name: userName})
	switch err {
	case controller.ErrTokenExpired:
		log.Infoln("token expired")
		err := refreshToken(h, w, req, userName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "cannot refresh token",
			})

			return err
		}

	case controller.ErrNotFound:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no user,password pair",
		})

		return err

	case nil:
		attachTokenToResponse(w, user.Token, h.config.CookieDomain, user.ExpiresUTCNano)

	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot get user",
		})

		return err
	}

	json.NewEncoder(w).Encode(servermodel.ServerResponse{
		Status: "ok",
		Description: "authenticated",
	})

	return nil
}

// TODO: Invalidate stale sessions.
// User logs in in one device, get token,
// then logs in in another device, receives new token.
// Old token invalidates, though not expired...
// Check if stage.lesnotes.space tokens influences on
// lesnotes.space (it has .lesnotes.space domain)
func (h *Handler) Auth(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	cookie, err := req.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no token",
		})

		return err
	}

	token := cookie.Value

	user, err := h.controller.GetUser(req.Context(), log, &model.GetUserParams{Token: token})
	if err == controller.ErrTokenExpired {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ServerAuthorizeResponse{
			ServerResponse: servermodel.ServerResponse{
				Status: "error",
				Description: "token expired",
			},
			Expired: true,
		})

		return err
	}

	if err != nil {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "user not found",
		})

		return err
	}

	json.NewEncoder(w).Encode(model.ServerAuthorizeResponse{
		ServerResponse: servermodel.ServerResponse{
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

	return nil
}

func (h *Handler) GetUser(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	values := req.URL.Query()
	if values.Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "empty user id",
		})

		return errors.New("empty user id")
	}

	id, err := strconv.Atoi(values.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "invalid id",
		})

		return err
	}

	user, err := h.controller.GetUser(req.Context(), log, &model.GetUserParams{ID: int32(id)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot find user",
		})

		return err
	}

	json.NewEncoder(w).Encode(model.ServerUserResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "exists",
		},
		User: model.User{
			ID: user.ID,
			Name: user.Name,
		},
	})

	return nil
}

func (h *Handler) Signup(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	userName, ok := getTextField(w, req, "name")
	if !ok {
		return errors.New("no user name")
	}

	exists, err := h.controller.HasUser(req.Context(), log, &model.HasUserParams{
		User: &model.User{Name: userName},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot check user",
		})

		return err
	}

	if exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "name exists",
		})

		return errors.New("name exists")
	}

	password, ok := getTextField(w, req, "password")
	if !ok {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no password",
		})

		return errors.New("cannot get password from request")
	}

	token, expiresUtcNano := createToken(w, h.config.CookieDomain)

	err = h.controller.AddUser(req.Context(), log, &model.AddUserParams{
		User: &model.User{
			Name:                  userName,
			Password:              password,
			Token:                 token,
			ExpiresUTCNano:        expiresUtcNano,
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot add user",
		})

		return err
	}

	json.NewEncoder(w).Encode(servermodel.ServerResponse{
		Status: "ok",
		Description: "created",
	})

	return nil
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

func refreshToken(h *Handler, w http.ResponseWriter, req *http.Request, userName string) error {
	token, expiresUtcNano := createToken(w, h.config.CookieDomain)

	return h.controller.RefreshToken(req.Context(), logger.Default(), &model.RefreshTokenParams{
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