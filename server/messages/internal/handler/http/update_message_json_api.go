package http

import (
	"net/http"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	users "github.com/bd878/gallery/server/users/pkg/model"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) UpdateMessageJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
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

	var request messages.UpdateRequest
	if err = json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to parse message",
			},
		})

		return
	}

	if request.MessageID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoID,
				Explain: "empty message id",
			},
		})

		return nil
	}

	var (
		threadID int64
		public   int
		text, title, name string
	)

	if request.ThreadID != nil {
		threadID = *request.ThreadID
	} else {
		threadID = -1
	}

	if request.Public != nil {
		public = *request.Public
	} else {
		public = -1
	}

	if request.Text != nil {
		text = *request.Text
	}

	if request.Title != nil {
		title = *request.Title
	}

	if request.Name != nil {
		name = *request.Name
	}

	return h.updateMessage(req.Context(), w, request.MessageID, user, text, title, name, threadID, request.FileIDs, public)
}
