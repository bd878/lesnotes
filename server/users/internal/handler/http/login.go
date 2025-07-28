package http

import (
	"time"
	"net/http"
	"errors"
	"encoding/json"

	"github.com/bd878/gallery/server/users/internal/controller"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) Login(w http.ResponseWriter, req *http.Request) error {
	name, ok := getTextField(w, req, "name")
	if !ok {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "no name",
		})

		return errors.New("no name field")
	}

	password, ok := getTextField(w, req, "password")
	if !ok {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "no password",
		})

		return errors.New("no password field")
	}

	session, err := h.controller.LoginUser(req.Context(), name, password)
	switch err {
	case controller.ErrUserNotFound:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "no user,password pair",
		})

		return err

	case controller.ErrWrongPassword:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "wrong password",
		})

		return err

	case nil:
		// attach session to response
		http.SetCookie(w, &http.Cookie{
			Name:          "token",
			Value:          session.Token,
			Domain:         h.config.CookieDomain,
			Expires:        time.Unix(0, session.ExpiresUTCNano),
			Path:          "/",
			HttpOnly:       true,
		})

	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "cannot get user",
		})

		return err
	}

	json.NewEncoder(w).Encode(servermodel.ServerResponse{
		Status:      "ok",
		Description: "authenticated",
	})

	return nil
}
