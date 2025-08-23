package http

import (
	"net/http"
	"fmt"
	"encoding/json"

	"github.com/bd878/gallery/server/utils"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) PublishMessageOrMessagesJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
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
				Code:    server.CodeNoBody,
				Explain: "body data required",
			},
		})

		return fmt.Errorf("no req data")
	}

	var request messages.PublishRequest
	if err := json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:   server.CodeWrongFormat,
				Explain: "failed to parse message",
			},
		})

		return err
	}

	if request.MessageID != nil {
		return h.publishMessage(w, req, user, *request.MessageID)
	} else if request.IDs != nil {
		return h.publishMessages(w, req, user, *request.IDs)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoID,
				Explain: "empty message id or batch ids",
			},
		})

		return fmt.Errorf("empty message id or batch ids")
	}
}