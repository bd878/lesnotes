package http

import (
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"

	"github.com/bd878/gallery/server/utils"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
)

func (h *Handler) PublishMessageOrMessages(w http.ResponseWriter, req *http.Request) error {
	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:   server.CodeNoUser,
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
				Code: messages.CodeMessagePublic,
				Explain: "can publish messages of a user only",
			},
		})

		return fmt.Errorf("can publish messages of a user only")
	}

	values := req.URL.Query()
	if values.Get("id") != "" {
		id, err := strconv.Atoi(values.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:   &server.ErrorCode{
					Code: server.CodeNoID,
					Explain: "invalid id",
				},
			})

			return err
		}

		return h.publishMessage(w, req, user, int64(id))
	} else if values.Get("ids") != "" {
		var ids []int64
		if err := json.Unmarshal([]byte(values.Get("ids")), &ids); err != nil {
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error:  &server.ErrorCode{
					Code: server.CodeWrongQuery,
					Explain: "wrong \"ids\" query field format",
				},
			})

			return err
		}

		return h.publishMessages(w, req, user, ids)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoID,
				Explain: "empty message id or batch ids",
			},
		})

		return fmt.Errorf("empty message id or batch ids")
	}
}

func (h *Handler) publishMessage(w http.ResponseWriter, req *http.Request, user *users.User, id int64) error {
	res, err := h.controller.PublishMessages(req.Context(), &messages.PublishMessagesParams{
		IDs: []int64{id},
		UserID: user.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code: messages.CodePublishFailed,
				Explain: "failed to publish a message",
			},
		})

		return err
	}

	response, err := json.Marshal(messages.PublishResponse{
		UpdateUTCNano: res.UpdateUTCNano,
		IDs: []int64{id},
		Description: "published",
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

func (h *Handler) publishMessages(w http.ResponseWriter, req *http.Request, user *users.User, ids []int64) error {
	res, err := h.controller.PublishMessages(req.Context(), &messages.PublishMessagesParams{
		IDs: ids,
		UserID: user.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    messages.CodePublishFailed,
				Explain: "failed to publish batch messages",
			},
		})

		return err
	}

	response, err := json.Marshal(messages.PublishResponse{
		UpdateUTCNano: res.UpdateUTCNano,
		IDs: ids,
		Description: "published",
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