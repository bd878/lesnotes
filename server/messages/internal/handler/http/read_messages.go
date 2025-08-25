package http

import (
	"fmt"
	"context"
	"net/http"
	"strconv"
	"encoding/json"

	"github.com/bd878/gallery/server/logger"
	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	files "github.com/bd878/gallery/server/files/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) ReadMessages(w http.ResponseWriter, req *http.Request) (err error) {
	var (
		limit, offset, order, public int
		threadID, messageID int64
	)

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

			return
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

				return
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

			return
		}
	} else {
		public = -1
	}

	if values.Has("id") {
		id, err := strconv.Atoi(values.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:   &server.ErrorCode{
					Code:    server.CodeNoID,
					Explain: "invalid message id",
				},
			})

			return err
		}

		messageID = int64(id)
	}

	if values.Get("thread") != "" {
		id, err := strconv.Atoi(values.Get("thread"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    server.CodeNoID,
					Explain: "invalid thread",
				},
			})

			return err
		}

		threadID = int64(id)
	}

	if values.Has("asc") {
		order, err = strconv.Atoi(values.Get("asc"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:  &server.ErrorCode{
					Code:    server.CodeWrongQuery,
					Explain: fmt.Sprintf("wrong \"%s\" query param", "asc"),
				},
			})

			return
		}
	}

	return h.readMessageOrMessages(req.Context(), w, user, limit, offset, public, messageID, threadID, order)
}

func (h *Handler) readMessageOrMessages(ctx context.Context, w http.ResponseWriter, user *users.User,
	limit int, offset int, public int, messageID int64, threadID int64, order int,
) (err error) {
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

		return
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
				Code:     server.CodeWrongQuery,
				Explain: "both public and id params are given",
			},
		})

		return
	}

	if messageID != 0 && threadID != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongQuery,
				Explain: "both id and thread params are given",
			},
		})

		return
	}

	if messageID != 0 && (limit > 0 || offset > 0) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code: server.CodeWrongQuery,
				Explain: "both id and limit/offset params given",
			},
		})

		return
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
	message, err := h.controller.ReadMessage(ctx, messageID, []int64{user.ID, users.PublicUserID})
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

	var list []*files.File
	for _, id := range message.FileIDs {
		file, err := h.filesGateway.ReadFile(ctx, message.UserID, id)
		if err != nil {
			logger.Errorw("failed to read file for a message", "user_id", message.UserID, "file_id", id, "message_id", message.ID)
			continue
		}

		list = append(list, &files.File{
			Name: file.Name,
			ID:   file.ID,
		})
	}
	message.Files = list

	response, err := json.Marshal(messages.ReadResponse{
		Messages:   []*messages.Message{message},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status: "ok",
		Response: json.RawMessage(response),
	})

	return
}

func (h *Handler) readThreadMessages(ctx context.Context, w http.ResponseWriter, user *users.User, threadID int64, limit, offset int32, ascending bool, private int32) (err error) {
	list, isLastPage, err := h.controller.ReadThreadMessages(ctx, user.ID, threadID, limit, offset, ascending, private)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:   &server.ErrorCode{
				Code:     server.CodeReadFailed,
				Explain: "failed to read thread messages",
			},
		})

		return err
	}

	fileIDs := make([]int64, 0)
	for _, message := range list {
		if message.FileIDs != nil {
			// TODO: fileIDs set
			fileIDs = append(fileIDs, message.FileIDs...)
		}
	}

	filesRes, err := h.filesGateway.ReadBatchFiles(ctx, fileIDs, user.ID)
	if err != nil {
		logger.Errorw("failed to read batch files", "user_id", user.ID, "error", err)
	} else {
		for _, message := range list {
			var list []*files.File
			for _, id := range message.FileIDs {
				file := filesRes[id]
				if file != nil {
					list = append(list, &files.File{
						ID: file.ID,
						Name: file.Name,
					})
				}
			}
			message.Files = list

			if message.UserID == users.PublicUserID {
				message.UserID = 0
			}
		}
	}

	response, err := json.Marshal(messages.ReadResponse{
		ThreadID:   &threadID,
		Messages:   list,
		IsLastPage: &isLastPage,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status: "ok",
		Response: json.RawMessage(response),
	})

	return
}

func (h *Handler) readMessages(ctx context.Context, w http.ResponseWriter, user *users.User, limit, offset int32, ascending bool, private int32) (err error) {
	list, isLastPage, err := h.controller.ReadMessages(ctx, user.ID, limit, offset, ascending, private)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to read all messages",
			},
		})

		return err
	}

	fileIDs := make([]int64, 0)
	for _, message := range list {
		if message.FileIDs != nil {
			// TODO: fileIDs set
			fileIDs = append(fileIDs, message.FileIDs...)
		}
	}

	filesRes, err := h.filesGateway.ReadBatchFiles(ctx, fileIDs, user.ID)
	if err != nil {
		logger.Errorw("failed to read batch files", "user_id", user.ID, "error", err)
	} else {
		for _, message := range list {
			var list []*files.File
			for _, id := range message.FileIDs {
				file := filesRes[id]
				if file != nil {
					list = append(list, &files.File{
						ID: file.ID,
						Name: file.Name,
					})
				}
			}
			message.Files = list

			if message.UserID == users.PublicUserID {
				message.UserID = 0
			}
		}
	}

	response, err := json.Marshal(messages.ReadResponse{
		Messages:   list,
		IsLastPage: &isLastPage,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status: "ok",
		Response: json.RawMessage(response),
	})

	return
}