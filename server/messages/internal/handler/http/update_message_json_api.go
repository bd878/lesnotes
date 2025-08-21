package http

import (
	"net/http"
	"fmt"
	"encoding/json"

	servermodel "github.com/bd878/gallery/server/pkg/model"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

func (h *Handler) UpdateMessageJsonAPI(w http.ResponseWriter, req *http.Request) error {
	var threadID int64
	var public int
	var text string

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

	var jsonRequest model.UpdateMessageJsonRequest
	if err := json.Unmarshal(data, &jsonRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to parse message",
		})

		return err
	}

	if jsonRequest.MessageID == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "empty message id",
		})

		return nil
	}

	if jsonRequest.ThreadID != nil {
		threadID = *jsonRequest.ThreadID
	} else {
		threadID = -1
	}

	if jsonRequest.Public != nil {
		public = *jsonRequest.Public
	} else {
		public = -1
	}

	if jsonRequest.Text != nil {
		text = *jsonRequest.Text
	}

	return h.updateMessage(req.Context(), w, *jsonRequest.MessageID, user, text, threadID, public)
}
