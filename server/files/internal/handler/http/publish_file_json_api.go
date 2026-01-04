package http

import (
	"context"
	"net/http"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	files "github.com/bd878/gallery/server/files/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) PublishFileJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
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

	if user.ID == users.PublicUserID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    files.CodeFilePublic,
				Explain: "can publish file of a user only",
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

	var request files.PublishFileRequest
	if err = json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to parse request",
			},
		})

		return
	}

	if request.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoID,
				Explain: "empty id",
			},
		})
		return
	}

	return h.publishFile(req.Context(), w, user.ID, request.ID)
}

func (h *Handler) publishFile(ctx context.Context, w http.ResponseWriter, userID, id int64) (err error) {
	err = h.controller.PublishFile(ctx, id, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    files.CodePublishFailed,
				Explain: "failed to publish file",
			},
		})

		return err
	}

	response, err := json.Marshal(files.PublishFileResponse{
		ID:            id,
		Description:   "public",
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