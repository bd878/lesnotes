package http

import (
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"

	servermodel "github.com/bd878/gallery/server/pkg/model"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/logger"
)

func (h *Handler) PublishMessageOrMessages(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
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
			Description: "can publish messages of a user only",
		})

		return fmt.Errorf("can publish messages of a user only")
	}

	values := req.URL.Query()
	if values.Get("id") != "" {
		return h.publishMessage(log, w, req, user, values.Get("id"))
	} else if values.Get("ids") != "" {
		return h.publishMessages(log, w, req, user, values.Get("ids"))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "empty message id or batch ids",
		})

		return fmt.Errorf("empty message id or batch ids")
	}
}

func (h *Handler) publishMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request, user *usermodel.User, idValue string) error {
	id, err := strconv.Atoi(idValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "invalid id",
		})

		return err
	}

	res, err := h.controller.PublishMessages(req.Context(), log, &model.PublishMessagesParams{
		IDs: []int32{int32(id)},
		UserID: user.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to publish a message",
		})

		return err
	}

	json.NewEncoder(w).Encode(model.PublishMessagesServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "published",
		},
		UpdateUTCNano: res.UpdateUTCNano,
		IDs: []int32{int32(id)},
	})

	return nil
}

func (h *Handler) publishMessages(log *logger.Logger, w http.ResponseWriter, req *http.Request, user *usermodel.User, idsValue string) error {
	var ids []int32
	if err := json.Unmarshal([]byte(idsValue), &ids); err != nil {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "wrong \"ids\" query field format",
		})

		return err
	}

	res, err := h.controller.PublishMessages(req.Context(), log, &model.PublishMessagesParams{
		IDs: ids,
		UserID: user.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to publish batch messages",
		})

		return err
	}

	json.NewEncoder(w).Encode(model.PublishMessagesServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "published",
		},
		UpdateUTCNano: res.UpdateUTCNano,
		IDs: ids,
	})

	return nil
}