package http

import (
	"net/http"
	"fmt"
	"encoding/json"

	"github.com/bd878/gallery/server/utils"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) ReadMessageOrMessagesJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
	var public int

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeNoUser,
				Explain: "user required",
			},
		})

		return fmt.Errorf("no user")
	}

	data, ok := utils.GetJsonRequestBody(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code: server.CodeNoBody,
				Explain: "body data required",
			},
		})

		return fmt.Errorf("no req data")
	}

	var request messages.ReadRequest
	if err := json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: messages.CodeNoMessage,
				Explain: "failed to parse message",
			},
		})

		return err
	}

	if request.Public == nil {
		public = -1
	} else {
		public = *request.Public
	}

	return h.readMessageOrMessages(req.Context(), w, user,
		request.Limit, request.Offset, public, request.MessageID, request.ThreadID, request.Asc)
}