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
		limit, offset, order int
		threadID, messageID int64
		ids []int64
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

	if values.Get("ids") != "" {
		if err := json.Unmarshal([]byte(values.Get("ids")), &ids); err != nil {
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:  &server.ErrorCode{
					Code:    server.CodeWrongQuery,
					Explain: "wrong \"ids\" query field format",
				},
			})

			return err
		}
	}

	if values.Get("user") != "" {
		id, err := strconv.Atoi(values.Get("user"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    server.CodeNoID,
					Explain: "invalid user",
				},
			})

			return err
		}

		// read public messages only
		return h.readMessageOrMessages(req.Context(), w, int64(id), limit, offset, messageID, threadID, order, true)
	} else if len(ids) > 0 {
		// read batch messages by given ids
		return h.readBatchMessages(req.Context(), w, user.ID, ids)
	} else {
		// read both public and private messages, 
		return h.readMessageOrMessages(req.Context(), w, user.ID, limit, offset, messageID, threadID, order, false)
	}
}

func (h *Handler) readBatchMessages(ctx context.Context, w http.ResponseWriter, userID int64, ids []int64) (err error) {
	// TODO: implement
	return nil
}

func (h *Handler) readMessageOrMessages(ctx context.Context, w http.ResponseWriter, userID int64,
	limit int, offset int, messageID int64, threadID int64, order int, publicOnly bool,
) (err error) {
	var ascending bool

	if userID == users.PublicUserID && messageID == 0 {
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
				Code:    server.CodeWrongQuery,
				Explain: "both id and limit/offset params given",
			},
		})

		return
	}

	if messageID != 0 {
		// read one message
		return h.readMessage(ctx, w, userID, messageID, publicOnly)
	} else if threadID != 0 {
		// read thread messages
		return h.readThreadMessages(ctx, w, userID, threadID, int32(limit), int32(offset), ascending, publicOnly)
	} else {
		// read all messages
		return h.readMessages(ctx, w, userID, int32(limit), int32(offset), ascending, publicOnly)
	}
}

func (h *Handler) readMessage(ctx context.Context, w http.ResponseWriter, userID int64, messageID int64, publicOnly bool) (err error) {
	message, err := h.controller.ReadMessage(ctx, messageID, []int64{userID, users.PublicUserID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    messages.CodeReadFailed,
				Explain: "failed to read a message",
			},
		})

		return err
	}

	if message.Private && publicOnly {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    messages.CodeReadFailed,
				Explain: "cannot read private message",
			},
		})

		return
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

	if message.UserID == users.PublicUserID {
		message.UserID = 0
	}

	response, err := json.Marshal(messages.ReadResponse{
		Messages:   []*messages.Message{message},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status:   "ok",
		Response: json.RawMessage(response),
	})

	return
}

func filterPublicMessages(list []*messages.Message) (filtered []*messages.Message) {
	filtered = make([]*messages.Message, 0, len(list))
	for _, message := range list {
		if !message.Private {
			filtered = append(filtered, message)
		}
	}
	return
}

func (h *Handler) readThreadMessages(ctx context.Context, w http.ResponseWriter, userID int64, threadID int64, limit, offset int32, ascending bool, publicOnly bool) (err error) {
	// TODO: read if thread is public

	list, isLastPage, err := h.controller.ReadThreadMessages(ctx, userID, threadID, limit, offset, ascending)
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

	// filter public messages only
	if publicOnly {
		list = filterPublicMessages(list)
	}

	fileIDs := make([]int64, 0)
	for _, message := range list {
		if message.FileIDs != nil {
			// TODO: fileIDs set
			fileIDs = append(fileIDs, message.FileIDs...)
		}
	}

	filesRes, err := h.filesGateway.ReadBatchFiles(ctx, fileIDs, userID)
	if err != nil {
		logger.Errorw("failed to read batch files", "user_id", userID, "error", err)
	} else {
		for _, message := range list {
			var list []*files.File
			for _, id := range message.FileIDs {
				file := filesRes[id]
				if file != nil {
					list = append(list, &files.File{
						ID:   file.ID,
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
		Status:   "ok",
		Response: json.RawMessage(response),
	})

	return
}

func (h *Handler) readMessages(ctx context.Context, w http.ResponseWriter, userID int64, limit, offset int32, ascending, publicOnly bool) (err error) {
	list, isLastPage, err := h.controller.ReadMessages(ctx, userID, limit, offset, ascending)
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

	if publicOnly {
		list = filterPublicMessages(list)
	}

	fileIDs := make([]int64, 0)
	for _, message := range list {
		if message.FileIDs != nil {
			// TODO: fileIDs set
			fileIDs = append(fileIDs, message.FileIDs...)
		}
	}

	filesRes, err := h.filesGateway.ReadBatchFiles(ctx, fileIDs, userID)
	if err != nil {
		logger.Errorw("failed to read batch files", "user_id", userID, "error", err)
	} else {
		for _, message := range list {
			var list []*files.File
			for _, id := range message.FileIDs {
				file := filesRes[id]
				if file != nil {
					list = append(list, &files.File{
						ID:   file.ID,
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
		Status:   "ok",
		Response: json.RawMessage(response),
	})

	return
}