package http

import (
	"time"
	"net/http"
	"encoding/json"

	servermodel "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) Logout(w http.ResponseWriter, req *http.Request) error {
	cookie, err := req.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "bad cookie",
		})

		return err
	}

	token := cookie.Value

	err = h.controller.LogoutUser(req.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "failed to delete token",
		})

		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:           "token",
		Value:          "",
		Domain:         h.config.CookieDomain,
		Expires:        time.Unix(0, 0),
		Path:           "/",
		HttpOnly:       true,
	})

	json.NewEncoder(w).Encode(servermodel.ServerResponse{
		Status:      "ok",
		Description: "logged out",
	})

	return nil
}
