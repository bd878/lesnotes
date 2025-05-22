package http

import (
	"net/http"
	"errors"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/users/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

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

	_, err = h.controller.AddUser(req.Context(), log, &model.AddUserParams{
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
