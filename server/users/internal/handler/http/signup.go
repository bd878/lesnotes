package http

import (
	"time"
	"net/http"
	"errors"
	"encoding/json"

	"github.com/bd878/gallery/server/utils"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) Signup( w http.ResponseWriter, req *http.Request) error {
	name, ok := getTextField(w, req, "name")
	if !ok {
		return errors.New("no user name")
	}

	password, ok := getTextField(w, req, "password")
	if !ok {
		return errors.New("cannot get password from request")
	}

	id := utils.RandomID()

	user, err := h.controller.CreateUser(req.Context(), int64(id), name, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "cannot add user",
		})

		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:          "token",
		Value:          user.Token,
		Domain:         h.config.CookieDomain,
		Expires:        time.Unix(0, user.ExpiresUTCNano),
		Path:          "/",
		HttpOnly:       true,
	})

	json.NewEncoder(w).Encode(servermodel.ServerResponse{
		Status: "ok",
		Description: "created",
	})

	return nil
}
