package http

import (
	"net/http"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	usersmodel "github.com/bd878/gallery/server/users/pkg/model"
	threadsmodel "github.com/bd878/gallery/server/threads/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) ReadThreadJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
	user, ok := req.Context().Value(middleware.UserContextKey{}).(*usersmodel.User)
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

	var request threadsmodel.ReadThreadRequest
	if err := json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    threadsmodel.CodeNoRequest,
				Explain: "failed to parse read thread request",
			},
		})

		return err
	}

	var userID int64
	if request.UserID != 0 {
		userID = request.UserID
	} else {
		userID = user.ID
	}

	thread, err := h.controller.ReadThread(req.Context(), request.ID, userID, request.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to read thread",
			},
		})

		return err
	}

	response, err := json.Marshal(threadsmodel.ReadThreadResponse{
		Threads:     []*threadsmodel.Thread{thread},
		Description: "ok",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status:    "ok",
		Response:  json.RawMessage(response),
	})

	return
}