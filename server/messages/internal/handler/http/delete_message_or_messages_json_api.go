package http

import (
	"net/http"
	"fmt"
	"encoding/json"

	servermodel "github.com/bd878/gallery/server/pkg/model"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/logger"
)

func (h *Handler) DeleteMessageOrMessagesJsonAPI(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "user required",
		})

		return fmt.Errorf("no user")
	}

	data, ok := utils.GetJsonRequestBody(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "body data required",
		})

		return fmt.Errorf("no req data")
	}

	var jsonRequest model.PrivateMessageOrMessagesJsonRequest
	if err := json.Unmarshal(data, &jsonRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to parse message",
		})

		return err
	}

	if jsonRequest.MessageID != nil {
		return h.deleteMessage(log, w, req, user, *jsonRequest.MessageID)
	} else if jsonRequest.IDs != nil {
		return h.deleteMessages(log, w, req, user, *jsonRequest.IDs)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "empty message id or batch ids",
		})

		return fmt.Errorf("empty message id or batch ids")
	}
}