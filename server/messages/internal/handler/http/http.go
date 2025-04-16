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

func (h *Handler) SendMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
	var (
		err error
		fileID, threadID int32
		private, hasFile bool
	)

	if err = req.ParseMultipartForm(1); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to parse form",
		})

		return
	}

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "user required",
		})

		return
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

			return
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

		return
	}

	if fileID != 0 && hasFile {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "should be either file_id or file, not both",
		})

		return
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

			return
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

			return
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
			log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "cannot read file",
			})

			return
		}

		fileName := filepath.Base(fh.Filename)

		fileResult, err := h.filesGateway.SaveFile(req.Context(), log, f, &model.SaveFileParams{
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

	resp, err := h.controller.SaveMessage(req.Context(), log, &model.Message{
		UserID: user.ID,
		ThreadID: threadID,
		Text:   text,
		FileID: fileID,
		Private: private,
	})
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to save message",
		})

		return
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
			log.Errorw("failed to read file for a message", "user_id", user.ID, "file_id", fileID, "message_id", resp.ID, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "failed to read a file",
			})
			return
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
}

func (h *Handler) DeleteMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) {
	user, ok := utils.GetUser(w, req)
	if !ok {
		log.Error("user not found")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "user required",
		})

		return
	}

	values := req.URL.Query()
	if values.Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "empty message id",
		})

		return
	}

	id, err := strconv.Atoi(values.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
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
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to delete message",
		})

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
	var private bool

	user, ok := utils.GetUser(w, req)
	if !ok {
		log.Error("user not found")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "user required",
		})

		return
	}

	values := req.URL.Query()
	if values.Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "empty message id",
		})

		return
	}

	id, err := strconv.Atoi(values.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "invalid id",
		})

		return
	}

	text := req.PostFormValue("text")
	if text == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "empty message field",
		})

		return
	}

	if values.Get("public") != "" {
		public, err := strconv.Atoi(values.Get("public"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "invalid public param",
			})

			return
		}

		if public > 0 {
			private = false
		} else if public == 0 {
			private = true
		} else {
			private = true
		}
	}

	resp, err := h.controller.UpdateMessage(req.Context(), log, &model.UpdateMessageParams{
		ID:     int32(id),
		UserID: user.ID,
		Text:   text,
		Private: private,
	})
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to update message",
		})

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
	var (
		limitInt, offsetInt, orderInt, publicInt int
		threadID, messageID, privateInt int32
		ascending bool
		message *model.Message
		messages []*model.Message
		isLastPage bool
		err error
	)

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
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
		}
	}

	if values.Has("public") {
		publicInt, err = strconv.Atoi(values.Get("public"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: fmt.Sprintf("wrong \"%s\" query param", "public"),
			})

			return
		}

		if publicInt > 0 {
			privateInt = 0
		} else {
			privateInt = 1
		}
	} else {
		privateInt = -1
	}

	if user.ID == usermodel.PublicUserID {
		privateInt = -1
	}

	if user.ID == usermodel.PublicUserID && !values.Has("id") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "can not list public messages",
		})

		return
	}

	if values.Get("id") != "" {
		messageid, err := strconv.Atoi(values.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
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
				Status: "error",
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

	if values.Has("public") && values.Has("id") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "both public and id params are given",
		})

		return
	}

	if messageID != 0 && threadID != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "both message_id and thread_id params are given",
		})

		return
	}

	if messageID != 0 && (limitInt > 0 || offsetInt > 0) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "both message_id and limit/offset params given",
		})

		return
	}

	log.Infow("read messages", "user_id", user.ID, "thread_id", threadID, "message_id", messageID, "public", publicInt)

	if messageID != 0 {
		// read one message
		res, err := h.controller.ReadOneMessage(req.Context(), log, &model.ReadOneMessageParams{
			ID: messageID,
			UserIDs: []int32{user.ID, usermodel.PublicUserID},
		})
		if err != nil {
			log.Errorw("failed to read one message", "user_id", user.ID, "message_id", messageID, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "failed to read a message",
			})

			return
		}

		message = res

		if user.ID == usermodel.PublicUserID {
			message.UserID = 0
		}

	} else if threadID != 0 {
		// read thread messages
		res, err := h.controller.ReadThreadMessages(req.Context(), log, &model.ReadThreadMessagesParams{
			UserID:    user.ID,
			ThreadID:  threadID,
			Limit:     int32(limitInt),
			Offset:    int32(offsetInt),
			Ascending: ascending,
			Private:   privateInt,
		})
		if err != nil {
			log.Errorw("failed to read thread messages, controller returned error", "user_id", user.ID, "thread_id", threadID, "public", publicInt, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "failed to read thread messages",
			})

			return
		}

		messages = res.Messages
		isLastPage = res.IsLastPage
	} else {
		// read all messages
		res, err := h.controller.ReadAllMessages(req.Context(), log, &model.ReadMessagesParams{
			UserID:    user.ID,
			Limit:     int32(limitInt),
			Offset:    int32(offsetInt),
			Ascending: ascending,
			Private:   privateInt,
		})
		if err != nil {
			log.Errorw("failed to read all messages", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "failed to read all messages",
			})

			return
		}

		messages = res.Messages
		isLastPage = res.IsLastPage
	}

	log.Infow("read messages", "user_id", user.ID, "name", user.Name, "message_id", messageID, "thread_id", threadID, "len_messages", len(messages), "public", publicInt)

	if message != nil {
		// one message
		if message.FileID != 0 {
			fileRes, err := h.filesGateway.ReadFile(req.Context(), log, user.ID, message.FileID)
			if err != nil {
				log.Errorw("failed to read file for a message", "user_id", user.ID, "file_id", message.FileID, "message_id", messageID, "error", err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(servermodel.ServerResponse{
					Status: "error",
					Description: "failed to read a file",
				})
				return
			}
			message.File = &filesmodel.File{
				Name: fileRes.Name,
				ID: message.FileID,
			}
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
		if message.FileID != 0 {
			fileIds = append(fileIds, message.FileID)
		}
	}

	filesRes, err := h.filesGateway.ReadBatchFiles(req.Context(), log, &model.ReadBatchFilesParams{
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
		if message.FileID != 0 {
			file := filesRes.Files[message.FileID]
			if file != nil {
				message.File = &filesmodel.File{
					ID: file.ID,
					Name: file.Name,
				}
			}
		}

		if user.ID == usermodel.PublicUserID {
			message.UserID = 0
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
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to get status",
		})

		return
	}
}
