package http

import (
	"net/http"
	"time"
	"encoding/json"

	"github.com/bd878/gallery/server/users/internal/controller"
	"github.com/bd878/gallery/server/users/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) Auth(w http.ResponseWriter, req *http.Request) error {
	cookie, err := req.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "no token",
		})

		return err
	}

	token := cookie.Value

	user, err := h.controller.AuthUser(req.Context(), token)
	if err == controller.ErrTokenExpired {
		w.WriteHeader(http.StatusUnauthorized)

		http.SetCookie(w, &http.Cookie{
			Name:           "token",
			Value:          "",
			Domain:         h.config.CookieDomain,
			Expires:        time.Unix(0, 0),
			Path:           "/",
			HttpOnly:       true,
		})

		json.NewEncoder(w).Encode(model.ServerAuthorizeResponse{
			ServerResponse: servermodel.ServerResponse{
				Status:      "error",
				Description: "token expired",
			},
			Expired: true,
		})

		return err
	}

	if err != nil {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "user not found",
		})

		return err
	}

	json.NewEncoder(w).Encode(model.ServerAuthorizeResponse{
		ServerResponse: servermodel.ServerResponse{
			Status:      "ok",
			Description: "token valid",
		},
		Expired: false,
		User: model.User{
			ID:               user.ID,
			Login:            user.Login,
			Token:            user.Token,
			ExpiresUTCNano:   user.ExpiresUTCNano,
		},
	})

	return nil
}
