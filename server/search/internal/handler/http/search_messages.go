package http

import (
	"fmt"
	"context"
	"net/http"
	"strconv"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	server "github.com/bd878/gallery/server/pkg/model"
	usersmodel "github.com/bd878/gallery/server/users/pkg/model"
	searchmodel "github.com/bd878/gallery/server/search/pkg/model"
)

func (h *Handler) SearchMessages(w http.ResponseWriter, req *http.Request) (err error) {
	var (
		query string
		threadID int64
		public int32
	)

	user, ok := req.Context().Value(middleware.UserContextKey{}).(*usersmodel.User)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoUser,
				Explain: "user required",
			},
		})

		return
	}

	values := req.URL.Query()

	if values.Has("query") {
		query = values.Get("query")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    searchmodel.CodeNoQuery,
				Explain: "no query",
			},
		})
		return
	}

	if values.Has("thread") {
		id, err := strconv.Atoi(values.Get("thread"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    server.CodeWrongQuery,
					Explain: fmt.Sprintf("wrong \"%s\" query param", "thread"),
				},
			})

			return err
		}

		threadID = int64(id)
	}

	if values.Has("public") {
		value, err := strconv.Atoi(values.Get("public"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    server.CodeWrongQuery,
					Explain: fmt.Sprintf("wrong \"%s\" query param", "public"),
				},
			})

			return err
		}

		public = int32(value)
	} else {
		public = -1
	}

	return h.searchMessages(req.Context(), w, user.ID, query, threadID, public)
}

func (h *Handler) searchMessages(ctx context.Context, w http.ResponseWriter, userID int64, query string, threadID int64, public int32) (err error) {
	list, err := h.controller.SearchMessages(ctx, userID, query, threadID, public)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to search messages",
			},
		})

		return err
	}

	response, err := json.Marshal(searchmodel.SearchMessagesResponse{
		Messages:    list,
		Count:       int32(len(list)),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	json.NewEncoder(w).Encode(server.ServerResponse{
		Status:    "ok",
		Response:  json.RawMessage(response),
	})

	return
}
