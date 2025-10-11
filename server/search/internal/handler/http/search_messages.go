package http

import (
	"context"
	"net/http"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	server "github.com/bd878/gallery/server/pkg/model"
	usersmodel "github.com/bd878/gallery/server/users/pkg/model"
	searchmodel "github.com/bd878/gallery/server/search/pkg/model"
)

func (h *Handler) SearchMessages(w http.ResponseWriter, req *http.Request) (err error) {
	var query string

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

	return h.searchMessages(req.Context(), w, user.ID, query)
}

func (h *Handler) searchMessages(ctx context.Context, w http.ResponseWriter, userID int64, query string) (err error) {
	list, err := h.controller.SearchMessages(ctx, userID, query)
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
