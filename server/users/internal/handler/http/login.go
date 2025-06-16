package http

import (
	"net/http"
	"errors"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/users/internal/controller"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

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
		_, err := refreshToken(h, w, req, userName)
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
