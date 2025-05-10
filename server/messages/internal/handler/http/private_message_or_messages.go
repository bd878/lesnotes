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

func (h *Handler) PrivateMessageOrMessages(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
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
			Description: "can private messages of a user only",
		})

		return fmt.Errorf("can private messages of a user only")
	}

	values := req.URL.Query()
	if values.Get("id") != "" {
		return h.privateMessage(log, w, req, user, values.Get("id"))
	} else if values.Get("ids") != "" {
		return h.privateMessages(log, w, req, user, values.Get("ids"))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "empty message id or batch ids",
		})

		return fmt.Errorf("empty message id or batch ids")
	}
}

func (h *Handler) privateMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request, user *usermodel.User, idValue string) error {
	id, err := strconv.Atoi(idValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "invalid id",
		})

		return err
	}

	res, err := h.controller.PrivateMessages(req.Context(), log, &model.PrivateMessagesParams{
		IDs: []int32{int32(id)},
		UserID: user.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to private a message",
		})

		return err
	}

	json.NewEncoder(w).Encode(model.PrivateMessagesServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "private",
		},
		UpdateUTCNano: res.UpdateUTCNano,
		IDs: []int32{int32(id)},
	})

	return nil
}

func (h *Handler) privateMessages(log *logger.Logger, w http.ResponseWriter, req *http.Request, user *usermodel.User, idsValue string) error {
	var ids []int32
	if err := json.Unmarshal([]byte(idsValue), &ids); err != nil {
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "wrong \"ids\" query field format",
		})

		return err
	}

	res, err := h.controller.PrivateMessages(req.Context(), log, &model.PrivateMessagesParams{
		IDs: ids,
		UserID: user.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to private batch messages",
		})

		return err
	}

	json.NewEncoder(w).Encode(model.PrivateMessagesServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "private",
		},
		UpdateUTCNano: res.UpdateUTCNano,
		IDs: ids,
	})

	return nil
}