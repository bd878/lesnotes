package http

import (
	"io"
	"errors"
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/files/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) DownloadFile(w http.ResponseWriter, req *http.Request) error {
	var (
		fileID int32
	)

	values := req.URL.Query()
	if values.Get("id") != "" {
		fileid, err := strconv.Atoi(values.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "id is empty",
			})
			return err
		}
		fileID = int32(fileid)
	}

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "user required",
		})
		return errors.New("user required")
	}

	file, stream, err := h.controller.ReadFileStream(req.Context(), &model.ReadFileStreamParams{FileID: fileID, UserID: user.ID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "failed to read file",
		})
		return err
	}

	logger.Infow("downloading file", "name", file.Name)

	w.Header().Set("Content-Disposition", "attachment; " + "filename*=UTF-8''" + file.Name)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))

	_, err = io.Copy(w, stream)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status:      "error",
			Description: "failed to write file to response",
		})
		return err
	}

	return nil
}