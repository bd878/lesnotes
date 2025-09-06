package http

import (
	"net/http"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/bd878/gallery/server/utils"
	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) SendMessageJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
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

	var request messages.SendRequest
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

	if request.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    messages.CodeNoText,
				Explain: "text required",
			},
		})

		return
	}

	if request.ThreadID != 0 {
		w.WriteHeader(http.StatusNotImplemented)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeNoID,
				Explain: "cannot save with thread_id yet",
			},
		})

		return
	}

	private := true
	if user.ID == users.PublicUserID {
		private = false
	}

	id := utils.RandomID()
	name := uuid.New().String()

	// TODO: check file by file_id exists
	// TODO: check thread by thread_id exists

	return h.saveMessage(w, req, int64(id), request.Text, request.Title, request.FileIDs, request.ThreadID, user.ID, private, name)
}
