package http

import (
	"fmt"
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
		limitInt, offsetInt, orderInt, publicInt int
		threadID, messageID, privateInt int32
		ascending, ok bool
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
		limitInt, err = strconv.Atoi(values.Get("limit"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: fmt.Sprintf("wrong \"%s\" query param", "limit"),
			})

			return err
		}

		if values.Has("offset") {
			offsetInt, err = strconv.Atoi(values.Get("offset"))
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
		publicInt, err = strconv.Atoi(values.Get("public"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: fmt.Sprintf("wrong \"%s\" query param", "public"),
			})

			return err
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

		return fmt.Errorf("cannot list public messages")
	}

	if values.Get("id") != "" {
		messageid, err := strconv.Atoi(values.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "invalid message id",
			})

			return err
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

			return err
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

			return err
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

	if messageID != 0 && (limitInt > 0 || offsetInt > 0) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "both message_id and limit/offset params given",
		})

		return fmt.Errorf("both message_id and limit/offset params given")
	}

	if messageID != 0 {
		// read one message
		return h.readMessage(log, w, req, user, messageID)
	} else if threadID != 0 {
		// read thread messages
		return h.readThreadMessages(log, w, req, user, threadID, int32(limitInt), int32(offsetInt), ascending, privateInt)
	} else {
		// read all messages
		return h.readMessages(log, w, req, user, int32(limitInt), int32(offsetInt), ascending, privateInt)
	}
}

func (h *Handler) readMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request, user *usermodel.User, messageID int32) error {
	message, err := h.controller.ReadOneMessage(req.Context(), log, &model.ReadOneMessageParams{
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
		fileRes, err := h.filesGateway.ReadFile(req.Context(), log, user.ID, message.FileID)
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

func (h *Handler) readThreadMessages(log *logger.Logger, w http.ResponseWriter, req *http.Request, user *usermodel.User, threadID, limit, offset int32, ascending bool, private int32) error {
	res, err := h.controller.ReadThreadMessages(req.Context(), log, &model.ReadThreadMessagesParams{
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

	filesRes, err := h.filesGateway.ReadBatchFiles(req.Context(), log, &model.ReadBatchFilesParams{
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

func (h *Handler) readMessages(log *logger.Logger, w http.ResponseWriter, req *http.Request, user *usermodel.User, limit, offset int32, ascending bool, private int32) error {
	res, err := h.controller.ReadAllMessages(req.Context(), log, &model.ReadMessagesParams{
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

	filesRes, err := h.filesGateway.ReadBatchFiles(req.Context(), log, &model.ReadBatchFilesParams{
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