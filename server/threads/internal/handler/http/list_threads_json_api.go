package http

import (
	"net/http"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	usersmodel "github.com/bd878/gallery/server/users/pkg/model"
	threadsmodel "github.com/bd878/gallery/server/threads/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) ListThreadsJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
	_, ok := req.Context().Value(middleware.UserContextKey{}).(*usersmodel.User)
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

	var request threadsmodel.ListThreadsRequest
	if err := json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    threadsmodel.CodeNoRequest,
				Explain: "failed to parse list threads request",
			},
		})

		return err
	}

	ids, isLastPage, err := h.controller.ListThreads(req.Context(), request.UserID, request.ParentID,
		request.Limit, request.Offset, request.Asc)
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

	response, err := json.Marshal(threadsmodel.ListThreadsResponse{
		IDs:        ids,
		IsLastPage: isLastPage,
		Count:      int32(len(ids)),
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