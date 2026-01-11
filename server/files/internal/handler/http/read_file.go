package http

import (
	"io"
	"fmt"
	"encoding/json"
	"net/http"

	users "github.com/bd878/gallery/server/users/pkg/model"
	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	files "github.com/bd878/gallery/server/files/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) ReadFile(w http.ResponseWriter, req *http.Request) (err error) {
	var (
		fileName    string
		isPublic bool
	)

	if req.PathValue("name") != "" {
		fileName = req.PathValue("name")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:     server.CodeWrongFormat,
				Explain:  "name param required",
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
	} else {
		isPublic = false
	}

	file, stream, err := h.controller.ReadFileStream(req.Context(), 0, fileName, isPublic)
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

	w.Header().Set("Content-Disposition", "inline; " + "filename*=UTF-8''" + file.Name)
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