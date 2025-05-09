package http

import (
	"net/http"
	"strconv"
	"io"
	"fmt"
	"context"
	"path/filepath"
	"encoding/json"

	filesmodel "github.com/bd878/gallery/server/files/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/logger"
)

type Controller interface {
	ReadOneMessage(ctx context.Context, log *logger.Logger, params *model.ReadOneMessageParams) (*model.Message, error)
	SaveMessage(ctx context.Context, log *logger.Logger, message *model.Message) (*model.SaveMessageResult, error)
	UpdateMessage(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) (*model.UpdateMessageResult, error)
	DeleteMessage(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) (*model.DeleteMessageResult, error)
	DeleteMessages(ctx context.Context, log *logger.Logger, params *model.DeleteMessagesParams) (*model.DeleteMessagesResult, error)
	ReadAllMessages(ctx context.Context, log *logger.Logger, params *model.ReadMessagesParams) (*model.ReadMessagesResult, error)
	ReadThreadMessages(ctx context.Context, log *logger.Logger, params *model.ReadThreadMessagesParams) (*model.ReadThreadMessagesResult, error)
}

type FilesGateway interface {
	ReadBatchFiles(ctx context.Context, log *logger.Logger, params *model.ReadBatchFilesParams) (*model.ReadBatchFilesResult, error)
	ReadFile(ctx context.Context, log *logger.Logger, userID, fileID int32) (*filesmodel.File, error)
	SaveFile(ctx context.Context, log *logger.Logger, stream io.Reader, params *model.SaveFileParams) (*model.SaveFileResult, error)
}

type Handler struct {
	controller    Controller
	filesGateway  FilesGateway
}

func New(controller Controller, filesGateway FilesGateway) *Handler {
	return &Handler{controller, filesGateway}
}

func (h *Handler) SendMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	var (
		err error
		fileID, threadID int32
		private, hasFile bool
	)

	if err = req.ParseMultipartForm(1); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to parse form",
		})

		return err
	}

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "user required",
		})

		return fmt.Errorf("no user")
	}

	text := req.PostFormValue("text")

	if req.PostFormValue("file_id") != "" {
		fileid, err := strconv.Atoi(req.PostFormValue("file_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "invalid file_id",
			})

			return err
		}

		fileID = int32(fileid)
	}

	if _, ok := req.MultipartForm.File["file"]; ok {
		hasFile = true
	}

	log.Infow("received message", "text", text, "name", user.Name, "file_id", fileID, "has_file", hasFile)

	if fileID == 0 && !hasFile && text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "text or file_id or file required",
		})

		return nil
	}

	if fileID != 0 && hasFile {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "should be either file_id or file, not both",
		})

		return nil
	}

	values := req.URL.Query()
	if values.Get("thread_id") != "" {
		threadid, err := strconv.Atoi(values.Get("thread_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "invalid thread id",
			})

			return err
		}

		threadID = int32(threadid)
	}

	if values.Has("public") {
		public, err := strconv.Atoi(values.Get("public"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "invalid public param",
			})

			return err
		}

		if public > 0 {
			private = false
		} else if public == 0 {
			private = true
		} else {
			private = true
		}
	} else {
		private = true
	}

	if user.ID == usermodel.PublicUserID {
		private = false
	}

	if hasFile {
		f, fh, err := req.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "cannot read file",
			})

			return err
		}

		fileName := filepath.Base(fh.Filename)

		fileResult, err := h.filesGateway.SaveFile(req.Context(), log, f, &model.SaveFileParams{
			UserID: user.ID,
			Name:   fileName,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "cannot save file",
			})

			return err
		}

		fileID = fileResult.ID
	}

	resp, err := h.controller.SaveMessage(req.Context(), log, &model.Message{
		UserID: user.ID,
		ThreadID: threadID,
		Text:   text,
		FileID: fileID,
		Private: private,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to save message",
		})

		return err
	}

	message := &model.Message{
		ID:            resp.ID,
		ThreadID:      threadID,
		FileID:        fileID,
		Text:          text,
		UpdateUTCNano: resp.UpdateUTCNano,
		CreateUTCNano: resp.CreateUTCNano,
		Private:       resp.Private,
	}

	if fileID != 0 {
		fileRes, err := h.filesGateway.ReadFile(req.Context(), log, user.ID, fileID)
		if err != nil {
			log.Errorw("failed to read file for a message", "user_id", user.ID, "file_id", fileID, "message_id", resp.ID)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "failed to read a file",
			})
			return err
		}
		message.File = &filesmodel.File{
			Name: fileRes.Name,
			ID: fileID,
		}
	}

	json.NewEncoder(w).Encode(model.NewMessageServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "accepted",
		},
		Message: message,
	})

	return nil
}

func (h *Handler) UpdateMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	var private int32

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "user required",
		})

		return fmt.Errorf("user not found")
	}

	values := req.URL.Query()
	if values.Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "empty message id",
		})

		return nil
	}

	id, err := strconv.Atoi(values.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "invalid id",
		})

		return err
	}

	text := req.PostFormValue("text")
	if text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "empty text field",
		})

		return nil
	}

	public := req.PostFormValue("public")
	if public != "" {
		publicInt, err := strconv.Atoi(public)
		if err != nil {
			log.Errorw("wrong public param", "public", public)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "invalid public param",
			})

			return err
		}

		if publicInt == 1 {
			private = 0
		} else if publicInt == 0 {
			private = 1
		} else {
			private = -1
		}		
	} else {
		private = -1
	}

	if user.ID == usermodel.PublicUserID && private == 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot make private public message",
		})

		return nil
	}

	resp, err := h.controller.UpdateMessage(req.Context(), log, &model.UpdateMessageParams{
		ID:     int32(id),
		UserID: user.ID,
		Text:   text,
		Private: private,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to update message",
		})

		return err
	}

	json.NewEncoder(w).Encode(model.UpdateMessageServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "accepted",
		},
		ID: resp.ID,
		UpdateUTCNano: resp.UpdateUTCNano,
		Private: resp.Private,
	})

	return nil
}

func (h *Handler) GetStatus(log *logger.Logger, w http.ResponseWriter, _ *http.Request) error {
	if _, err := io.WriteString(w, "ok\n"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to get status",
		})

		return err
	}

	return nil
}
