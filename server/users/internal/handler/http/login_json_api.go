package http

import (
	"net/http"
	"io"
	"errors"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/users/internal/controller"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) LoginJsonAPI(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	var err error

	data, err := io.ReadAll(req.Body)
	if err != nil {
		req.Body.Close()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to parse request",
		})

		return err
	}

	defer req.Body.Close()

	var user model.User
	if err := json.Unmarshal(data, &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to parse user",
		})

		return err
	}

	if user.Name == "" {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no name",
		})

		return errors.New("cannot get user name from request")
	}

	if user.Password == "" {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "no password",
		})

		return errors.New("cannot get password from request")
	}


	exists, err := h.controller.HasUser(req.Context(), log, &model.HasUserParams{
		User: &user,
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

	resp, err := h.controller.GetUser(req.Context(), log, &model.GetUserParams{Name: user.Name})

	user.ExpiresUTCNano = resp.ExpiresUTCNano
	user.Token = resp.Token
	user.ID = resp.ID

	switch err {
	case controller.ErrTokenExpired:
		log.Infoln("token expired")
		err := refreshToken(h, w, req, user.Name)
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