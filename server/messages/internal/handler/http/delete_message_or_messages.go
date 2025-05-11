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
	"github.com/bd878/gallery/server/logger"
)

func (h *Handler) DeleteMessageOrMessages(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
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
		return h.deleteMessage(log, w, req, user, values.Get("id"))
	} else if values.Get("ids") != "" {
		return h.deleteMessages(log, w, req, user, values.Get("ids"))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "empty message id or batch ids",
		})

		return fmt.Errorf("empty message id or batch ids")
	}
}

func (h *Handler) deleteMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request, user *usermodel.User, idValue string) error {
	id, err := strconv.Atoi(idValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "invalid id",
		})

		return err
	}

	_, err = h.controller.DeleteMessage(req.Context(), log, &model.DeleteMessageParams{
		ID: int32(id),
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
		ID: int32(id),
	})

	return nil
}

func (h *Handler) deleteMessages(log *logger.Logger, w http.ResponseWriter, req *http.Request, user *usermodel.User, idsValue string) error {
	var ids []int32
	if err := json.Unmarshal([]byte(idsValue), &ids); err != nil {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "wrong \"ids\" query field format",
		})

		return err
	}

	res, err := h.controller.DeleteMessages(req.Context(), log, &model.DeleteMessagesParams{
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