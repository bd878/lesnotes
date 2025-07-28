package http

import (
	"net/http"
	"fmt"
	"encoding/json"

	servermodel "github.com/bd878/gallery/server/pkg/model"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/utils"
)

func (h *Handler) ReadMessageOrMessagesJsonAPI(w http.ResponseWriter, req *http.Request) error {
	var public int

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

	var jsonRequest model.ReadMessageOrMessagesJsonRequest
	if err := json.Unmarshal(data, &jsonRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to parse message",
		})

		return err
	}

	if jsonRequest.Public == nil {
		public = -1
	} else {
		public = *jsonRequest.Public
	}

	return h.readMessageOrMessages(req.Context(), w, user,
		jsonRequest.Limit, jsonRequest.Offset, public, jsonRequest.MessageID, jsonRequest.ThreadID, jsonRequest.Asc)
}