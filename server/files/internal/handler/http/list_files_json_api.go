package http

import (
	"context"
	"net/http"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	users "github.com/bd878/gallery/server/users/pkg/model"
	files "github.com/bd878/gallery/server/files/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) ListFilesJsonAPI(w http.ResponseWriter, req *http.Request) (err error) {
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

	var request files.ListFilesRequest
	if err := json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    files.CodeNoRequest,
				Explain: "failed to parse request",
			},
		})

		return err
	}

	if user.ID != users.PublicUserID {
		if request.Asc == 1 {
			return h.listPrivateFiles(req.Context(), w, user.ID, request.Limit, request.Offset, true)
		}
		return h.listPrivateFiles(req.Context(), w, user.ID, request.Limit, request.Offset, true)
	}

	if user.ID == users.PublicUserID && (request.UserID == 0 || request.UserID == users.PublicUserID) {
		if request.Asc == 1 {
			return h.listPublicFiles(req.Context(), w, request.Limit, request.Offset, true)
		}
		return h.listPublicFiles(req.Context(), w, request.Limit, request.Offset, false)
	}

	if user.ID == users.PublicUserID && request.UserID != 0 {
		if request.Asc == 1 {
			return h.listPublicUserFiles(req.Context(), w, request.UserID, request.Limit, request.Offset, true)
		}
		return h.listPublicUserFiles(req.Context(), w, request.UserID, request.Limit, request.Offset, false)
	}

	// empty
	response, err := json.Marshal(files.ListFilesResponse{
		Files:       make([]*files.File, 0),
		Description: "empty",
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status:     "ok",
		Response:   json.RawMessage(response),
	})

	return
}

func (h *Handler) listPrivateFiles(ctx context.Context, w http.ResponseWriter, userID int64, limit, offset int32, ascending bool) (err error) {
	list, err := h.controller.ListFiles(ctx, userID, limit, offset, ascending, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to read all files",
			},
		})

		return err
	}

	response, err := json.Marshal(files.ListFilesResponse{
		Files:       list.Files,
		IsLastPage:  list.IsLastPage,
		IsFirstPage: list.IsFirstPage,
		Offset:      list.Offset,
		Total:       list.Total,
		Count:       list.Count,
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

func (h *Handler) listPublicUserFiles(ctx context.Context, w http.ResponseWriter, userID int64, limit, offset int32, ascending bool) (err error) {
	list, err := h.controller.ListFiles(ctx, userID, limit, offset, ascending, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to read user public files",
			},
		})

		return err
	}

	response, err := json.Marshal(files.ListFilesResponse{
		Files:       list.Files,
		IsLastPage:  list.IsLastPage,
		IsFirstPage: list.IsFirstPage,
		Offset:      list.Offset,
		Total:       list.Total,
		Count:       list.Count,
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

func (h *Handler) listPublicFiles(ctx context.Context, w http.ResponseWriter, limit, offset int32, ascending bool) (err error) {
	list, err := h.controller.ListFiles(ctx, users.PublicUserID, limit, offset, ascending, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to read all public files",
			},
		})

		return err
	}

	response, err := json.Marshal(files.ListFilesResponse{
		Files:       list.Files,
		IsLastPage:  list.IsLastPage,
		IsFirstPage: list.IsFirstPage,
		Total:       list.Total,
		Offset:      list.Offset,
		Count:       list.Count,
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
