package http

import (
	"net/http"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	users "github.com/bd878/gallery/server/users/pkg/model"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) DeleteMessageOrMessagesJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
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

	data, ok := req.Context().Value(middleware.RequestContextKey{}).(json.RawMessage)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeNoBody,
				Explain: "request required",
			},
		})

		return
	}

	var request messages.DeleteRequest
	if err = json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to parse request",
			},
		})

		return
	}

	if request.IDs != nil {
		return h.deleteMessages(w, req, user, request.IDs)
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(server.ServerResponse{
		Status: "error",
		Error:  &server.ErrorCode{
			Code:    server.CodeWrongFormat,
			Explain: "empty ids",
		},
	})

	return
}