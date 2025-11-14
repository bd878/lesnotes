package http

import (
	"net/http"
	"strconv"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	threads "github.com/bd878/gallery/server/threads/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) PrivateThread(w http.ResponseWriter, req *http.Request) (err error) {
	user, ok := req.Context().Value(middleware.UserContextKey{}).(*users.User)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoUser,
				Explain: "user required",
			},
		})

		return
	}

	if user.ID == users.PublicUserID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    threads.CodeThreadPublic,
				Explain: "can private thread of a user only",
			},
		})

		return
	}

	values := req.URL.Query()
	if values.Get("id") != "" {
		id, err := strconv.Atoi(values.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:  &server.ErrorCode{
					Code:    server.CodeWrongFormat,
					Explain: "wrong id",
				},
			})

			return err
		}

		return h.privateThread(w, req, user, int64(id))
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(server.ServerResponse{
		Status: "error",
		Error: &server.ErrorCode{
			Code:    server.CodeWrongFormat,
			Explain: "no id",
		},
	})

	return
}

func (h *Handler) privateThread(w http.ResponseWriter, req *http.Request, user *users.User, id int64) (err error) {
	err = h.controller.PrivateThread(req.Context(), id, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    threads.CodePrivateFailed,
				Explain: "failed to private thread",
			},
		})

		return err
	}

	response, err := json.Marshal(threads.PrivateResponse{
		ID:            id,
		Description:   "private",
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