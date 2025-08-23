package http

import (
	"fmt"
	"context"
	"net/http"
	"strconv"
	"encoding/json"

	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/logger"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	files "github.com/bd878/gallery/server/files/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) ReadMessageOrMessages(w http.ResponseWriter, req *http.Request) (err error) {
	var (
		limit, offset, order, public int
		threadID, messageID int64
	)

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeNoUser,
				Explain: "user required",
			},
		})

		return fmt.Errorf("user required")
	}

	values := req.URL.Query()

	if values.Has("limit") {
		limit, err = strconv.Atoi(values.Get("limit"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    server.CodeWrongQuery,
					Explain: fmt.Sprintf("wrong \"%s\" query param", "limit"),
				},
			})

			return err
		}

		if values.Has("offset") {
			offset, err = strconv.Atoi(values.Get("offset"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(server.ServerResponse{
					Status: "error",
					Error: &server.ErrorCode{
						Code:    server.CodeWrongQuery,
						Explain: fmt.Sprintf("wrong \"%s\" query param", "offset"),
					},
				})

				return err
			}
		}
	}

	if values.Has("public") {
		public, err = strconv.Atoi(values.Get("public"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:  &server.ErrorCode{
					Code:    server.CodeWrongQuery,
					Explain: fmt.Sprintf("wrong \"%s\" query param", "public"),
				},
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
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:   &server.ErrorCode{
					Code: server.CodeNoID,
					Explain: "invalid message id",
				},
			})

			return err
		}

		messageID = int64(messageIDInt)
	}

	if values.Get("thread_id") != "" {
		threadid, err := strconv.Atoi(values.Get("thread_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code: server.CodeNoID,
					Explain: "invalid thread id",
				},
			})

			return err
		}

		threadID = int64(threadid)
	}

	if values.Has("asc") {
		order, err = strconv.Atoi(values.Get("asc"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:  &server.ErrorCode{
					Code: server.CodeWrongQuery,
					Explain: fmt.Sprintf("wrong \"%s\" query param", "asc"),
				},
			})

			return err
		}
	}

	return h.readMessageOrMessages(req.Context(), w, user, limit, offset, public, messageID, threadID, order)
}

func (h *Handler) readMessageOrMessages(ctx context.Context, w http.ResponseWriter, user *users.User,
	limit int, offset int, public int, messageID int64, threadID int64, order int,
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

	if user.ID == users.PublicUserID {
		private = -1
	}

	if user.ID == users.PublicUserID && messageID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    messages.CodeMessagePublic,
				Explain: "can not list public messages",
			},
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
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:   &server.ErrorCode{
				Code:        server.CodeWrongQuery,
				Explain: "both public and id params are given",
			},
		})

		return fmt.Errorf("both public and id params are given")
	}

	if messageID != 0 && threadID != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongQuery,
				Explain: "both message_id and thread_id params are given",
			},
		})

		return fmt.Errorf("both message_id and thread_id params are given")
	}

	if messageID != 0 && (limit > 0 || offset > 0) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code: server.CodeWrongQuery,
				Explain: "both message_id and limit/offset params given",
			},
		})

		return fmt.Errorf("both message_id and limit/offset params given")
	}

	if messageID != 0 {
		// read one message
		return h.readMessage(ctx, w, user, messageID)
	} else if threadID != 0 {
		// read thread messages
		return h.readThreadMessages(ctx, w, user, threadID, int32(limit), int32(offset), ascending, int32(private))
	} else {
		// read all messages
		return h.readMessages(ctx, w, user, int32(limit), int32(offset), ascending, int32(private))
	}
}

func (h *Handler) readMessage(ctx context.Context, w http.ResponseWriter, user *users.User, messageID int64) (err error) {
	message, err := h.controller.ReadOneMessage(ctx, &messages.ReadOneMessageParams{
		ID: messageID,
		UserIDs: []int64{user.ID, users.PublicUserID},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code: messages.CodeReadFailed,
				Explain: "failed to read a message",
			},
		})

		return err
	}

	if message.UserID == users.PublicUserID {
		message.UserID = 0
	}

	if message.FileID != 0 {
		fileRes, err := h.filesGateway.ReadFile(ctx, user.ID, message.FileID)
		if err != nil {
			logger.Errorw("failed to read file for a message", "user_id", user.ID, "file_id", message.FileID, "message_id", messageID, "error", err)
		} else {
			message.File = &files.File{
				Name: fileRes.Name,
				ID: message.FileID,
			}
		}
	}

	response, err := json.Marshal(messages.ReadResponse{
		Message:   message,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status: "ok",
		Response: json.RawMessage(response),
	})

	return nil
}

func (h *Handler) readThreadMessages(ctx context.Context, w http.ResponseWriter, user *users.User, threadID int64, limit, offset int32, ascending bool, private int32) error {
	res, err := h.controller.ReadThreadMessages(ctx, &messages.ReadThreadMessagesParams{
		UserID:    user.ID,
		ThreadID:  threadID,
		Limit:     limit,
		Offset:    offset,
		Ascending: ascending,
		Private:   private,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:   &server.ErrorCode{
				Code: messages.CodeReadFailed,
				Explain: "failed to read thread messages",
			},
		})

		return err
	}

	isLastPage := res.IsLastPage

	fileIds := make([]int64, 0)
	for _, message := range res.Messages {
		if message.FileID != 0 {
			fileIds = append(fileIds, message.FileID)
		}
	}

	filesRes, err := h.filesGateway.ReadBatchFiles(ctx, &messages.ReadBatchFilesParams{
		UserID: user.ID,
		IDs:    fileIds,
	})
	if err != nil {
		logger.Errorw("failed to read batch files", "user_id", user.ID, "error", err)
	} else {
		for _, message := range res.Messages {
			if message.FileID != 0 {
				file := filesRes.Files[message.FileID]
				if file != nil {
					message.File = &files.File{
						ID: file.ID,
						Name: file.Name,
					}
				}
			}

			if message.UserID == users.PublicUserID {
				message.UserID = 0
			}
		}
	}

	response, err := json.Marshal(messages.ListResponse{
		ThreadID: &threadID,
		Messages: res.Messages,
		IsLastPage: isLastPage,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status: "ok",
		Response: json.RawMessage(response),
	})

	return nil
}

func (h *Handler) readMessages(ctx context.Context, w http.ResponseWriter, user *users.User, limit, offset int32, ascending bool, private int32) error {
	res, err := h.controller.ReadAllMessages(ctx, &messages.ReadMessagesParams{
		UserID:    user.ID,
		Limit:     limit,
		Offset:    offset,
		Ascending: ascending,
		Private:   private,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code: server.CodeWrongFormat,
				Explain: "failed to read all messages",
			},
		})

		return err
	}

	isLastPage := res.IsLastPage

	fileIds := make([]int64, 0)
	for _, message := range res.Messages {
		if message.FileID != 0 {
			fileIds = append(fileIds, message.FileID)
		}
	}

	filesRes, err := h.filesGateway.ReadBatchFiles(ctx, &messages.ReadBatchFilesParams{
		UserID: user.ID,
		IDs:    fileIds,
	})
	if err != nil {
		logger.Errorw("failed to read batch files", "user_id", user.ID, "error", err)
	} else {
		for _, message := range res.Messages {
			if message.FileID != 0 {
				file := filesRes.Files[message.FileID]
				if file != nil {
					message.File = &files.File{
						ID: file.ID,
						Name: file.Name,
					}
				}
			}

			if message.UserID == users.PublicUserID {
				message.UserID = 0
			}
		}
	}

	response, err := json.Marshal(messages.ListResponse{
		Messages: res.Messages,
		IsLastPage: isLastPage,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status: "ok",
		Response: json.RawMessage(response),
	})

	return nil
}