package http

import (
	"net/http"
	"fmt"
	"encoding/json"

	"github.com/bd878/gallery/server/utils"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) SendMessageJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:   server.CodeNoUser,
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
				Code:  server.CodeNoBody,
				Explain: "body data required",
			},
		})

		return fmt.Errorf("no req data")
	}

	var message messages.Message
	if err = json.Unmarshal(data, &message); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:  server.CodeWrongFormat,
				Explain: "failed to parse message",
			},
		})

		return err
	}

	if message.Text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    messages.CodeNoText,
				Explain: "text required",
			},
		})

		return nil
	}

	if message.ThreadID != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeNoID,
				Explain: "cannot save with thread_id yet",
			},
		})

		return nil
	}

	if user.ID == users.PublicUserID {
		message.Private = false
	}

	message.UserID = user.ID

	// TODO: check file by file_id exists
	// TODO: check thread by thread_id exists

	return h.saveMessage(w, req, &message)
}
