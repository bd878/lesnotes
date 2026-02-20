package http

import (
	"io"
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/internal/logger"
	users "github.com/bd878/gallery/server/users/pkg/model"
	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	files "github.com/bd878/gallery/server/files/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) DownloadFile(w http.ResponseWriter, req *http.Request) (err error) {
	var (
		fileID int64
		isPublic bool
	)

	values := req.URL.Query()

	if values.Has("id") {
		value, err := strconv.Atoi(values.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    server.CodeWrongQuery,
					Explain: "wrong id param",
				},
			})
			return err
		}
		fileID = int64(value)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeNoID,
				Explain: "id is empty",
			},
		})
		return
	}

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
		isPublic = true
	}

	file, stream, err := h.controller.ReadFileStream(req.Context(), fileID, "", isPublic)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:    files.CodeReadFailed,
				Explain: "failed to read file",
			},
		})
		return err
	}

	logger.Infow("downloading file", "name", file.Name)

	w.Header().Set("Content-Disposition", "attachment; " + "filename*=UTF-8''" + file.Name)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))

	_, err = io.Copy(w, stream)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status:      "error",
			Error: &server.ErrorCode{
				Code:    files.CodeWriteFailed,
				Explain: "failed to write file to response",
			},
		})
		return
	}

	return
}