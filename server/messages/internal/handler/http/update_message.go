package http

import (
	"net/http"
	"context"
	"strconv"
	"fmt"
	"encoding/json"

	"github.com/bd878/gallery/server/utils"
	server "github.com/bd878/gallery/server/pkg/model"
	users "github.com/bd878/gallery/server/users/pkg/model"
	messages "github.com/bd878/gallery/server/messages/pkg/model"
)

func (h *Handler) UpdateMessage(w http.ResponseWriter, req *http.Request) (err error) {
	var public, threadID int

	user, ok := utils.GetUser(w, req)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoUser,
				Explain: "user required",
			},
		})

		return fmt.Errorf("user not found")
	}

	values := req.URL.Query()
	if values.Get("id") == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoID,
				Explain: "empty message id",
			},
		})

		return nil
	}

	id, err := strconv.Atoi(values.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: server.CodeNoID,
				Explain: "invalid id",
			},
		})

		return err
	}

	text := req.PostFormValue("text")

	if req.PostFormValue("thread_id") != "" {
		threadID, err = strconv.Atoi(req.PostFormValue("thread_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code: server.CodeNoID,
					Explain: "invalid thread id",
				},
			})

			return err
		}
	} else {
		threadID = -1
	}

	if req.PostFormValue("public") != "" {
		public, err = strconv.Atoi(req.PostFormValue("public"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code: server.CodeWrongFormat,
					Explain: "invalid public param",
				},
			})

			return err
		}
	} else {
		public = -1
	}

	return h.updateMessage(req.Context(), w, int64(id), user, text, int64(threadID), public)
}

func (h *Handler) updateMessage(ctx context.Context, w http.ResponseWriter, messageID int64, user *users.User,
	text string, threadID int64, public int,
) (err error) {
	var private int32

	if public == 1 {
		private = 0
	} else if public == 0 {
		private = 1
	} else {
		private = -1
	}

	if user.ID == users.PublicUserID && threadID != -1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: messages.CodeMessagePublic,
				Explain: "cannot move public message",
			},
		})

		return
	}

	if user.ID == users.PublicUserID && private == 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code: messages.CodeMessagePublic,
				Explain: "cannot make private public message",
			},
		})

		return
	}

	msg, err := h.controller.ReadOneMessage(ctx, &messages.ReadOneMessageParams{
		ID: messageID,
		UserIDs: []int64{user.ID},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code: messages.CodeNoMessage,
				Explain: "failed to read message",
			},
		})

		return err		
	}

	if text == "" && threadID == -1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "text or thread_id or both must be provided",
			},
		})

		return nil
	}

	if text == "" {
		text = msg.Text
	}

	resp, err := h.controller.UpdateMessage(ctx, &messages.UpdateMessageParams{
		ID:     messageID,
		UserID: user.ID,
		Text:   text,
		Private: private,
		ThreadID: threadID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    messages.CodeUpdateFailed,
				Explain: "failed to update message",
			},
		})

		return err
	}

	response, err := json.Marshal(messages.UpdateResponse{
		Description: "updated",
		ID: resp.ID,
		UpdateUTCNano: resp.UpdateUTCNano,
		Private: resp.Private,
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