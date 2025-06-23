package http

import (
	"fmt"
	"context"
	"net/http"
	"strconv"
	"encoding/json"

	filesmodel "github.com/bd878/gallery/server/files/pkg/model"
	servermodel "github.com/bd878/gallery/server/pkg/model"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/logger"
)

func (h *Handler) ReadMessageOrMessages(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	var (
		limit, offset, order, public int
		threadID, messageID int32
		err error
	)

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "user required",
		})

		return fmt.Errorf("user required")
	}

	values := req.URL.Query()

	if values.Has("limit") {
		limit, err = strconv.Atoi(values.Get("limit"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: fmt.Sprintf("wrong \"%s\" query param", "limit"),
			})

			return err
		}

		if values.Has("offset") {
			offset, err = strconv.Atoi(values.Get("offset"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(servermodel.ServerResponse{
					Status: "error",
					Description: fmt.Sprintf("wrong \"%s\" query param", "offset"),
				})

				return err
			}
		}
	}

	if values.Has("public") {
		public, err = strconv.Atoi(values.Get("public"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: fmt.Sprintf("wrong \"%s\" query param", "public"),
			})

			return err
		}
	} else {
		public = -1
	}

	if values.Has("id") {
		messageIDInt, err := strconv.Atoi(values.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "invalid message id",
			})

			return err
		}

		messageID = int32(messageIDInt)
	}

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

	if values.Has("asc") {
		order, err = strconv.Atoi(values.Get("asc"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: fmt.Sprintf("wrong \"%s\" query param", "asc"),
			})

			return err
		}
	}

	return h.readMessageOrMessages(req.Context(), log, w, user, limit, offset, public, messageID, threadID, order)
}

func (h *Handler) readMessageOrMessages(ctx context.Context, log *logger.Logger, w http.ResponseWriter, user *usermodel.User,
	limit int, offset int, public int, messageID int32, threadID int32, order int,
) error {
	var (
		private int
		ascending bool
	)

	if public > 0 {
		private = 0
	} else if public == 0 {
		private = 1
	} else {
		private = -1
	}

	if user.ID == usermodel.PublicUserID {
		private = -1
	}

	if user.ID == usermodel.PublicUserID && messageID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "can not list public messages",
		})

		return fmt.Errorf("cannot list public messages")
	}

	switch order {
	case 0:
		ascending = false
	case 1:
		ascending = true
	default:
		ascending = true
	}

	if public > 0 && messageID > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "both public and id params are given",
		})

		return fmt.Errorf("both public and id params are given")
	}

	if messageID != 0 && threadID != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "both message_id and thread_id params are given",
		})

		return fmt.Errorf("both message_id and thread_id params are given")
	}

	if messageID != 0 && (limit > 0 || offset > 0) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "both message_id and limit/offset params given",
		})

		return fmt.Errorf("both message_id and limit/offset params given")
	}

	if messageID != 0 {
		// read one message
		return h.readMessage(ctx, log, w, user, messageID)
	} else if threadID != 0 {
		// read thread messages
		return h.readThreadMessages(ctx, log, w, user, threadID, int32(limit), int32(offset), ascending, int32(private))
	} else {
		// read all messages
		return h.readMessages(ctx, log, w, user, int32(limit), int32(offset), ascending, int32(private))
	}
}

func (h *Handler) readMessage(ctx context.Context, log *logger.Logger, w http.ResponseWriter, user *usermodel.User, messageID int32) error {
	message, err := h.controller.ReadOneMessage(ctx, log, &model.ReadOneMessageParams{
		ID: messageID,
		UserIDs: []int32{user.ID, usermodel.PublicUserID},
	})
	if err != nil {
		log.Errorw("failed to read one message", "user_id", user.ID, "message_id", messageID)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to read a message",
		})

		return err
	}

	if message.UserID == usermodel.PublicUserID {
		message.UserID = 0
	}

	if message.FileID != 0 {
		fileRes, err := h.filesGateway.ReadFile(ctx, log, user.ID, message.FileID)
		if err != nil {
			log.Errorw("failed to read file for a message", "user_id", user.ID, "file_id", message.FileID, "message_id", messageID, "error", err)
		} else {
			message.File = &filesmodel.File{
				Name: fileRes.Name,
				ID: message.FileID,
			}
		}
	}

	json.NewEncoder(w).Encode(model.ReadMessageServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
		},
		Message: message,
	})

	return nil
}

func (h *Handler) readThreadMessages(ctx context.Context, log *logger.Logger, w http.ResponseWriter, user *usermodel.User, threadID, limit, offset int32, ascending bool, private int32) error {
	res, err := h.controller.ReadThreadMessages(ctx, log, &model.ReadThreadMessagesParams{
		UserID:    user.ID,
		ThreadID:  threadID,
		Limit:     limit,
		Offset:    offset,
		Ascending: ascending,
		Private:   private,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to read thread messages",
		})

		return err
	}

	messages := res.Messages
	isLastPage := res.IsLastPage

	fileIds := make([]int32, 0)
	for _, message := range messages {
		if message.FileID != 0 {
			fileIds = append(fileIds, message.FileID)
		}
	}

	filesRes, err := h.filesGateway.ReadBatchFiles(ctx, log, &model.ReadBatchFilesParams{
		UserID: user.ID,
		IDs:    fileIds,
	})
	if err != nil {
		log.Errorw("failed to read batch files", "user_id", user.ID, "error", err)
	} else {
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

			if message.UserID == usermodel.PublicUserID {
				message.UserID = 0
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

	return nil
}

func (h *Handler) readMessages(ctx context.Context, log *logger.Logger, w http.ResponseWriter, user *usermodel.User, limit, offset int32, ascending bool, private int32) error {
	res, err := h.controller.ReadAllMessages(ctx, log, &model.ReadMessagesParams{
		UserID:    user.ID,
		Limit:     limit,
		Offset:    offset,
		Ascending: ascending,
		Private:   private,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to read all messages",
		})

		return err
	}

	messages := res.Messages
	isLastPage := res.IsLastPage

	fileIds := make([]int32, 0)
	for _, message := range messages {
		if message.FileID != 0 {
			fileIds = append(fileIds, message.FileID)
		}
	}

	filesRes, err := h.filesGateway.ReadBatchFiles(ctx, log, &model.ReadBatchFilesParams{
		UserID: user.ID,
		IDs:    fileIds,
	})
	if err != nil {
		log.Errorw("failed to read batch files", "user_id", user.ID, "error", err)
	} else {
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

			if message.UserID == usermodel.PublicUserID {
				message.UserID = 0
			}
		}
	}

	json.NewEncoder(w).Encode(model.MessagesListServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
		},
		Messages: messages,
		IsLastPage: isLastPage,
	})

	return nil
}