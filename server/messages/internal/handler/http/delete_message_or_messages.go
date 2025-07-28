package http

import (
	"net/http"
	"fmt"
	"strconv"
	"encoding/json"

	servermodel "github.com/bd878/gallery/server/pkg/model"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/utils"
)

func (h *Handler) DeleteMessageOrMessages(w http.ResponseWriter, req *http.Request) error {
	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "user required",
		})

		return fmt.Errorf("user not found")
	}

	if user.ID == usermodel.PublicUserID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot delete public message",
		})

		return fmt.Errorf("cannot delete public message")
	}

	values := req.URL.Query()
	if values.Get("id") != "" {
		id, err := strconv.Atoi(values.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "invalid id",
			})

			return err
		}

		return h.deleteMessage(w, req, user, int32(id))
	} else if values.Get("ids") != "" {
		var ids []int32
		if err := json.Unmarshal([]byte(values.Get("ids")), &ids); err != nil {
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "wrong \"ids\" query field format",
			})

			return err
		}

		return h.deleteMessages(w, req, user, ids)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "empty message id or batch ids",
		})

		return fmt.Errorf("empty message id or batch ids")
	}
}

func (h *Handler) deleteMessage(w http.ResponseWriter, req *http.Request, user *usermodel.User, id int32) error {
	_, err := h.controller.DeleteMessage(req.Context(), &model.DeleteMessageParams{
		ID: id,
		UserID: user.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to delete message",
		})

		return err
	}

	json.NewEncoder(w).Encode(model.DeleteMessageServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "deleted",
		},
		ID: id,
	})

	return nil
}

func (h *Handler) deleteMessages(w http.ResponseWriter, req *http.Request, user *usermodel.User, ids []int32) error {
	res, err := h.controller.DeleteMessages(req.Context(), &model.DeleteMessagesParams{
		IDs: ids,
		UserID: user.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to delete batch messages",
		})

		return err
	}

	json.NewEncoder(w).Encode(model.DeleteMessagesServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "deleted",
		},
		IDs: res.IDs,
	})

	return nil
}