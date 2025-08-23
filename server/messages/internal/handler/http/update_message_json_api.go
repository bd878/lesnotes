package http

import (
	"net/http"
	"fmt"
	"errors"
	"encoding/json"

	"github.com/bd878/gallery/server/utils"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) UpdateMessageJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoUser,
				Explain: "user required",
			},
		})

		return errors.New("no user")
	}

	data, ok := utils.GetJsonRequestBody(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeWrongFormat,
				Explain: "body data required",
			},
		})

		return fmt.Errorf("no req data")
	}

	var request messages.UpdateRequest
	if err = json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeWrongFormat,
				Explain: "failed to parse message",
			},
		})

		return err
	}

	if request.MessageID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoID,
				Explain: "empty message id",
			},
		})

		return nil
	}

	var (
		threadID int64
		public   int
		text     string
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

	return h.updateMessage(req.Context(), w, request.MessageID, user, text, threadID, public)
}
