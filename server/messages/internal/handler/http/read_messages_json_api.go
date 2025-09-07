package http

import (
	"net/http"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	users "github.com/bd878/gallery/server/users/pkg/model"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) ReadMessagesJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
	var threadID int64

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

	var request messages.ReadRequest
	if err := json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    messages.CodeNoMessage,
				Explain: "failed to parse message",
			},
		})

		return err
	}

	if request.ThreadID != nil {
		threadID = *request.ThreadID
	} else {
		threadID = -1
	}

	if request.UserID != 0 {
		return h.readMessageOrMessages(req.Context(), w, request.UserID, request.Limit, request.Offset, request.MessageID, threadID, request.Asc, true)
	} else if len(request.IDs) > 0 {
		return h.readBatchMessages(req.Context(), w, user.ID, request.IDs)
	} else {
		return h.readMessageOrMessages(req.Context(), w, user.ID, request.Limit, request.Offset, request.MessageID, threadID, request.Asc, false)
	}
}