package http

import (
	"net/http"
	"fmt"
	"strconv"
	"encoding/json"

	"github.com/bd878/gallery/server/utils"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) DeleteMessageOrMessages(w http.ResponseWriter, req *http.Request) (err error) {
	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:  server.CodeNoUser,
				Explain: "user required",
			},
		})

		return fmt.Errorf("user not found")
	}

	if user.ID == users.PublicUserID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:  messages.CodeMessagePublic,
				Explain: "cannot delete public message",
			},
		})

		return fmt.Errorf("cannot delete public message")
	}

	values := req.URL.Query()
	if values.Get("id") != "" {
		id, err := strconv.Atoi(values.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:   server.CodeNoID,
					Explain: "invalid id",
				},
			})

			return err
		}

		return h.deleteMessage(w, req, user, int64(id))
	} else if values.Get("ids") != "" {
		var ids []int64
		if err := json.Unmarshal([]byte(values.Get("ids")), &ids); err != nil {
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:  &server.ErrorCode{
					Code:  server.CodeWrongQuery,
					Explain: "wrong \"ids\" query field format",
				},
			})

			return err
		}

		return h.deleteMessages(w, req, user, ids)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:   &server.ErrorCode{
				Code:  server.CodeWrongFormat,
				Explain: "empty message id or batch ids",
			},
		})

		return fmt.Errorf("empty message id or batch ids")
	}
}

func (h *Handler) deleteMessage(w http.ResponseWriter, req *http.Request, user *users.User, id int64) (err error) {
	_, err = h.controller.DeleteMessage(req.Context(), &messages.DeleteMessageParams{
		ID: id,
		UserID: user.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:  messages.CodeDeleteFailed,
				Explain: "failed to delete message",
			},
		})

		return
	}

	response, err := json.Marshal(messages.DeleteResponse{
		ID: &id,
		Description: "deleted",
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

func (h *Handler) deleteMessages(w http.ResponseWriter, req *http.Request, user *users.User, ids []int64) (err error) {
	_, err = h.controller.DeleteMessages(req.Context(), &messages.DeleteMessagesParams{
		IDs: ids,
		UserID: user.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code: messages.CodeDeleteFailed,
				Explain: "failed to delete batch messages",
			},
		})

		return err
	}

	response, err := json.Marshal(messages.DeleteResponse{
		IDs: &ids,
		Description: "deleted",
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