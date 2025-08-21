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
)

func (h *Handler) PrivateMessageOrMessages(w http.ResponseWriter, req *http.Request) error {
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
		id, err := strconv.Atoi(values.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "invalid id",
			})

			return err
		}

		return h.privateMessage(w, req, user, int64(id))
	} else if values.Get("ids") != "" {
		var ids []int64
		if err := json.Unmarshal([]byte(values.Get("ids")), &ids); err != nil {
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "wrong \"ids\" query field format",
			})

			return err
		}

		return h.privateMessages(w, req, user, ids)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "empty message id or batch ids",
		})

		return fmt.Errorf("empty message id or batch ids")
	}
}

func (h *Handler) privateMessage(w http.ResponseWriter, req *http.Request, user *usermodel.User, id int64) error {
	res, err := h.controller.PrivateMessages(req.Context(), &model.PrivateMessagesParams{
		IDs: []int64{id},
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
		IDs: []int64{id},
	})

	return nil
}

func (h *Handler) privateMessages(w http.ResponseWriter, req *http.Request, user *usermodel.User, ids []int64) error {
	res, err := h.controller.PrivateMessages(req.Context(), &model.PrivateMessagesParams{
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