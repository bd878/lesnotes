package http

import (
	"net/http"
	"strconv"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	users "github.com/bd878/gallery/server/users/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) Update(w http.ResponseWriter, req *http.Request) (err error) {
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

	newLogin, newTheme, newLang, newFontSize := req.PostFormValue("login"), req.PostFormValue("theme"), req.PostFormValue("language"), req.PostFormValue("font_size")

	fontSize, err := strconv.Atoi(newFontSize)
	if err != nil {
		// TODO: log error
		fontSize = 0
	}

	err = h.controller.UpdateUser(req.Context(), user.ID, newLogin, newTheme, newLang, int32(fontSize))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:   &server.ErrorCode{
				Code:    users.CodeUpdateFailed,
				Explain: "cannot update user",
			},
		})

		return err
	}

	response, err := json.Marshal(users.UpdateResponse{
		Description: "updated",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status:       "ok",
		Response:     json.RawMessage(response),
	})

	return nil
}
