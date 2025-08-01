package http

import (
	"net/http"
	"fmt"
	"encoding/json"

	servermodel "github.com/bd878/gallery/server/pkg/model"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

func (h *Handler) SendMessageJsonAPI(w http.ResponseWriter, req *http.Request) error {
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

	var message model.Message
	if err := json.Unmarshal(data, &message); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to parse message",
		})

		return err
	}

	if message.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "text required",
		})

		return nil
	}

	if message.ThreadID != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot save with thread_id yet",
		})

		return nil
	}

	if user.ID == usermodel.PublicUserID {
		message.Private = false
	}

	message.UserID = user.ID

	// TODO: check file by file_id exists
	// TODO: check thread by thread_id exists

	return h.saveMessage(w, req, &message)
}
