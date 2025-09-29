package http

import (
	"net/http"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	users "github.com/bd878/gallery/server/users/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) GetMe(w http.ResponseWriter, req *http.Request) (err error) {
	user, ok := req.Context().Value(middleware.UserContextKey{}).(*users.User)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoUser,
				Explain: "cannot find user",
			},
		})

		return
	}

	if user.ID == users.PublicUserID {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:   "error",
			Error:    &server.ErrorCode{
				Code:     server.CodeNoUser,
				Explain:  "not authorized",
			},
		})

		return
	}

	response, err := json.Marshal(users.GetMeResponse{
		ID:       user.ID,
		Login:    user.Login,
		Theme:    user.Theme,
		Lang:     user.Lang,
		FontSize: user.FontSize,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status:   "ok",
		Response: json.RawMessage(response),
	})

	return
}
