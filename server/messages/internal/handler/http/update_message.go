package http

import (
	"net/http"
	"strconv"
	"fmt"
	"encoding/json"

	servermodel "github.com/bd878/gallery/server/pkg/model"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	"github.com/bd878/gallery/server/messages/pkg/model"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/logger"
)

func (h *Handler) UpdateMessage(log *logger.Logger, w http.ResponseWriter, req *http.Request) error {
	var private, threadID int32

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "user required",
		})

		return fmt.Errorf("user not found")
	}

	values := req.URL.Query()
	if values.Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "empty message id",
		})

		return nil
	}

	id, err := strconv.Atoi(values.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "invalid id",
		})

		return err
	}

	text := req.PostFormValue("text")

	if req.PostFormValue("thread_id") != "" {
		threadid, err := strconv.Atoi(req.PostFormValue("thread_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "invalid thread id",
			})

			return err
		}

		threadID = int32(threadid)
	} else {
		threadID = -1
	}

	if user.ID == usermodel.PublicUserID && threadID != -1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot move public message",
		})

		return nil
	}

	public := req.PostFormValue("public")
	if public != "" {
		publicInt, err := strconv.Atoi(public)
		if err != nil {
			log.Errorw("wrong public param", "public", public)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(servermodel.ServerResponse{
				Status: "error",
				Description: "invalid public param",
			})

			return err
		}

		if publicInt == 1 {
			private = 0
		} else if publicInt == 0 {
			private = 1
		} else {
			private = -1
		}		
	} else {
		private = -1
	}

	if user.ID == usermodel.PublicUserID && private == 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "cannot make private public message",
		})

		return nil
	}

	msg, err := h.controller.ReadOneMessage(req.Context(), log, &model.ReadOneMessageParams{
		ID: int32(id),
		UserIDs: []int32{user.ID},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to read message",
		})

		return err		
	}

	if text == "" && threadID == -1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "text or thread_id or both must be provided",
		})

		return nil
	}

	if text == "" {
		text = msg.Text
	}

	resp, err := h.controller.UpdateMessage(req.Context(), log, &model.UpdateMessageParams{
		ID:     int32(id),
		UserID: user.ID,
		Text:   text,
		Private: private,
		ThreadID: threadID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(servermodel.ServerResponse{
			Status: "error",
			Description: "failed to update message",
		})

		return err
	}

	json.NewEncoder(w).Encode(model.UpdateMessageServerResponse{
		ServerResponse: servermodel.ServerResponse{
			Status: "ok",
			Description: "accepted",
		},
		ID: resp.ID,
		UpdateUTCNano: resp.UpdateUTCNano,
		Private: resp.Private,
	})

	return nil
}