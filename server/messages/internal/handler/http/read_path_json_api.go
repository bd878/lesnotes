package http

import (
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/internal/logger"
	files "github.com/bd878/gallery/server/files/pkg/model"
	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	users "github.com/bd878/gallery/server/users/pkg/model"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) ReadPathJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
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

	var request messages.ReadPathRequest
	if err := json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    messages.CodeNoMessage,
				Explain: "failed to parse request",
			},
		})

		return err
	}

	// TODO: add "thread" to list, front expects it in loadStack
	list, parentID, err := h.controller.ReadPath(req.Context(), user.ID, request.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to read messages path",
			},
		})

		return err
	}

	fileIDs := make([]int64, 0)
	for _, message := range list {
		if message.FileIDs != nil {
			// TODO: fileIDs set
			fileIDs = append(fileIDs, message.FileIDs...)
		}
	}

	filesRes, err := h.filesGateway.ReadBatchFiles(req.Context(), fileIDs, user.ID)
	if err != nil {
		logger.Errorw("failed to read batch files", "user_id", user.ID, "error", err)
	} else {
		for _, message := range list {
			var list []*files.File
			for _, id := range message.FileIDs {
				file := filesRes[id]
				if file != nil {
					list = append(list, &files.File{
						ID:   file.ID,
						Name: file.Name,
					})
				}
			}
			message.Files = list

			if message.UserID == users.PublicUserID {
				message.UserID = 0
			}
		}
	}

	response, err := json.Marshal(messages.ReadPathResponse{
		Messages:   list,
		ThreadID:   parentID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status:   "ok",
		Response: json.RawMessage(response),
	})

	return
}