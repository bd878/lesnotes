package http

import (
	"net/http"
	"time"
	"encoding/json"

	"github.com/bd878/gallery/server/users/internal/controller"
	users "github.com/bd878/gallery/server/users/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) Auth(w http.ResponseWriter, req *http.Request) (err error) {
	cookie, err := req.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoToken,
				Explain: "no token",
			},
		})

		return err
	}

	token := cookie.Value

	_, err = h.controller.AuthUser(req.Context(), token)
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

		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeTokenExpired,
				Explain: "token expired",
			},
		})

		return
	}

	if err != nil {
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoUser,
				Explain: "user not found",
			},
		})

		return
	}

	response, err := json.Marshal(users.AuthResponse{
		Expired: false,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status: "ok",
		Response: json.RawMessage(response),
	})

	return nil
}
