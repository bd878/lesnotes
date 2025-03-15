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
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/logger"
)

const selectNoLimit int = 0

type Controller interface {
	ReadOneMessage(ctx context.Context, log *logger.Logger, userID, messageID int32) (*model.Message, error)
	SaveMessage(ctx context.Context, log *logger.Logger, message *model.Message) (*model.SaveMessageResult, error)
	UpdateMessage(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) (*model.UpdateMessageResult, error)
	DeleteMessage(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) (*model.DeleteMessageResult, error)
	ReadAllMessages(ctx context.Context, log *logger.Logger, params *model.ReadAllMessagesParams) (*model.ReadAllMessagesResult, error)
	ReadThreadMessages(ctx context.Context, log *logger.Logger, params *model.ReadThreadMessagesParams) (*model.ReadThreadMessagesResult, error)
}

// TODO: move files into controller?
type FilesGateway interface {
	SaveFileStream(ctx context.Context, log *logger.Logger, stream io.Reader, params *model.SaveFileParams) (*model.SaveFileResult, error)
	ReadBatchFiles(ctx context.Context, log *logger.Logger, params *model.ReadBatchFilesParams) (*model.ReadBatchFilesResult, error)
	ReadFile(ctx context.Context, log *logger.Logger, userID, fileID int32) (*filesmodel.File, error)
	// TODO: DeleteFile
}

type Handler struct {
	controller    Controller
	filesGateway  FilesGateway
}

func New(controller Controller, filesGateway FilesGateway) *Handler {
	return &Handler{controller, filesGateway}
}

// TODO: rewrite on builder pattern
func (h *Handler) SendMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
	var err error
	if err = req.ParseMultipartForm(1); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "ok",
			Description: "user required",
		})
		return
	}

	var fileName string
	var fileID int32
	if _, ok := req.MultipartForm.File["file"]; ok {
		f, fh, err := req.FormFile("file")
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "cannot read file",
			})
			return
		}

		fileName = filepath.Base(fh.Filename)

		// TODO: create fileID here, pipe in goroutine
		fileResult, err := h.filesGateway.SaveFileStream(context.Background(), log, f, &model.SaveFileParams{
			UserID: user.ID,
			Name:   fileName,
		})
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "cannot save file",
			})
			return
		}

		fileID = fileResult.ID
	}

	text := req.PostFormValue("text")

	log.Infoln("text=", text, "user name=", user.Name)

	if fileName == "" && text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "text or file are empty",
		})
		return
	}

	var threadID int32
	values := req.URL.Query()
	if values.Get("thread_id") != "" {
		threadid, err := strconv.Atoi(values.Get("thread_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "ok",
				Description: "invalid thread id",
			})
			return
		}
		threadID = int32(threadid)
	}

	resp, err := h.controller.SaveMessage(req.Context(), log, &model.Message{
		UserID: user.ID,
		ThreadID: threadID,
		Text:   text,
		File:   &filesmodel.File{
			ID:     fileID,
		},
	})
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(model.NewMessageServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "accepted",
		},
		Message: &model.Message{
			ID:            resp.ID,
			ThreadID:      threadID,
			File:          &filesmodel.File{
				ID:          fileID,
				Name:        fileName,
			},
			Text:          text,
			UpdateUTCNano: resp.UpdateUTCNano,
			CreateUTCNano: resp.CreateUTCNano,
		},
	})
}

func (h *Handler) DeleteMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
	user, ok := utils.GetUser(w, req)
	if !ok {
		log.Error("user not found")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "ok",
			Description: "user required",
		})
		return
	}

	values := req.URL.Query()
	if values.Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "ok",
			Description: "empty message id",
		})
		return
	}

	id, err := strconv.Atoi(values.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "ok",
			Description: "invalid id",
		})
		return
	}

	_, err = h.controller.DeleteMessage(req.Context(), log, &model.DeleteMessageParams{
		ID: int32(id),
		UserID: user.ID,
	})
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(model.DeleteMessageServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "deleted",
		},
		ID: int32(id),
	})
}

func (h *Handler) UpdateMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
	user, ok := utils.GetUser(w, req)
	if !ok {
		log.Error("user not found")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "ok",
			Description: "user required",
		})
		return
	}

	values := req.URL.Query()
	if values.Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "ok",
			Description: "empty message id",
		})
		return
	}

	id, err := strconv.Atoi(values.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "ok",
			Description: "invalid id",
		})
		return
	}

	text := req.PostFormValue("text")
	if text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "ok",
			Description: "empty message field",
		})
		return
	}

	resp, err := h.controller.UpdateMessage(req.Context(), log, &model.UpdateMessageParams{
		ID:     int32(id),
		UserID: user.ID,
		Text:   text,
	})
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(model.UpdateMessageServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "accepted",
		},
		ID: resp.ID,
		UpdateUTCNano: resp.UpdateUTCNano,
	})
}

func (h *Handler) ReadMessagesOrMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
	var limitInt, offsetInt, orderInt int
	var threadID, messageID int32
	var ascending bool
	var message *model.Message
	var messages []*model.Message
	var isLastPage bool
	var err error

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "ok",
			Description: "user required",
		})
		return
	}

	values := req.URL.Query()

	if values.Has("limit") {
		limitInt, err = strconv.Atoi(values.Get("limit"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: fmt.Sprintf("wrong \"%s\" query param", "limit"),
			})
			return
		}

		if values.Has("offset") {
			offsetInt, err = strconv.Atoi(values.Get("offset"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(servermodel.ServerResponse{
					Status: "error",
					Description: fmt.Sprintf("wrong \"%s\" query param", "offset"),
				})
				return
			}
		} else {
			offsetInt = 0
		}
	} else {
		limitInt = selectNoLimit
	}

	if values.Get("message_id") != "" {
		messageid, err := strconv.Atoi(values.Get("message_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "ok",
				Description: "invalid message id",
			})
			return
		}
		messageID = int32(messageid)
	}

	if values.Get("thread_id") != "" {
		threadid, err := strconv.Atoi(values.Get("thread_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "ok",
				Description: "invalid thread id",
			})
			return
		}
		threadID = int32(threadid)
	}

	if values.Has("asc") {
		orderInt, err = strconv.Atoi(values.Get("asc"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: fmt.Sprintf("wrong \"%s\" query param", "asc"),
			})
			return
		}

		switch orderInt {
		case 0:
			ascending = false
		case 1:
			ascending = true
		default:
			ascending = true
		}
	} else {
		ascending = true
	}

	if messageID != 0 && threadID != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "both message_id and thread_id params given",
		})
		return
	}

	if messageID != 0 && (ascending != false || limitInt > 0 || offsetInt > 0) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "both message_id and limit/offset params given",
		})
		return
	}

	log.Infow("read messages", "user_id", user.ID, "thread_id", threadID)

	if threadID != 0 {
		// read thread messages
		res, err := h.controller.ReadThreadMessages(context.Background(), log, &model.ReadThreadMessagesParams{
			UserID:    user.ID,
			ThreadID:  threadID,
			Limit:     int32(limitInt),
			Offset:    int32(offsetInt),
			Ascending: ascending,
		})
		if err != nil {
			log.Errorw("failed to read thread messages, controller returned error", "user_id", user.ID, "thread_id", threadID, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		messages = res.Messages
		isLastPage = res.IsLastPage
	} else if messageID != 0 {
		// read one message
		res, err := h.controller.ReadOneMessage(context.Background(), log, user.ID, messageID)
		if err != nil {
			log.Errorw("failed to read one message", "user_id", user.ID, "message_id", messageID, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		message = res
	} else {
		// read all messages
		res, err := h.controller.ReadAllMessages(context.Background(), log, &model.ReadAllMessagesParams{
			UserID:    user.ID,
			Limit:     int32(limitInt),
			Offset:    int32(offsetInt),
			Ascending: ascending,
		})
		if err != nil {
			log.Errorw("failed to read all messages", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		messages = res.Messages
		isLastPage = res.IsLastPage
	}

	log.Infow("read messages", "thread_id", threadID, "len_messages", len(messages))

	if message != nil {
		// one message
		if message.File != nil && message.File.ID != 0 {
			fileRes, err := h.filesGateway.ReadFile(context.Background(), log, user.ID, message.File.ID)
			if err != nil {
				log.Errorw("failed to read file for a message", "user_id", user.ID, "message_id", messageID, "error", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(servermodel.ServerResponse{
					Status: "error",
					Description: "failed to read a file",
				})
				return
			}
			message.File.Name = fileRes.Name
		}

		json.NewEncoder(w).Encode(model.ReadMessageServerResponse{
			ServerResponse: servermodel.ServerResponse{
				Status: "ok",
			},
			Message: message,
		})
		return
	}


	fileIds := make([]int32, 0)
	for _, message := range messages {
		if message.File != nil && message.File.ID != 0 {
			fileIds = append(fileIds, message.File.ID)
		}
	}

	filesRes, err := h.filesGateway.ReadBatchFiles(context.Background(), log, &model.ReadBatchFilesParams{
		UserID: user.ID,
		IDs:    fileIds,
	})
	if err != nil {
		log.Errorw("failed to read batch files", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to read files",
		})
		return
	}

	for _, message := range messages {
		if message.File != nil {
			file := filesRes.Files[message.File.ID]
			if file != nil {
				message.File.Name = file.Name
			}
		}
	}

	json.NewEncoder(w).Encode(model.MessagesListServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
		},
		ThreadID: threadID,
		Messages: messages,
		IsLastPage: isLastPage,
	})
}

func (h *Handler) GetStatus(log *logger.Logger, w http.ResponseWriter, _ *http.Request) {
	if _, err := io.WriteString(w, "ok\n"); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
