package http

import (
	"time"
	"net/http"
	"encoding/json"

	users "github.com/bd878/gallery/server/users/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) Logout(w http.ResponseWriter, req *http.Request) (err error) {
	cookie, err := req.Cookie("token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code: users.CodeBadCookie,
				Explain: "bad cookie",
			},
		})

		return err
	}

	token := cookie.Value

	err = h.controller.LogoutUser(req.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:     users.CodeLogoutFailed,
				Explain: "failed to delete token",
			},
		})

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:           "token",
		Value:          "",
		Domain:         h.config.CookieDomain,
		Expires:        time.Unix(0, 0),
		Path:           "/",
		HttpOnly:       true,
	})

	response, err := json.Marshal(users.LogoutResponse{
		Description:    "ok",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status:      "ok",
		Response:    json.RawMessage(response),
	})

	return nil
}
