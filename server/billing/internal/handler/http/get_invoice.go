package http

import (
	"fmt"
	"net/http"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	usersmodel "github.com/bd878/gallery/server/users/pkg/model"
	billingmodel "github.com/bd878/gallery/server/billing/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) GetInvoice(w http.ResponseWriter, req *http.Request) (err error) {
	var invoiceID string

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

	if values.Has("id") {
		invoiceID = values.Get("id")
		if invoiceID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(server.ServerResponse{
				Status: "error",
				Error: &server.ErrorCode{
					Code:    server.CodeWrongQuery,
					Explain: fmt.Sprintf("wrong \"%s\" query param", "id"),
				},
			})

			return err
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoID,
				Explain: "no invoice id",
			},
		})

		return err
	}

	invoice, err := h.controller.GetInvoice(req.Context(), invoiceID, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to get invoice",
			},
		})

		return err
	}

	response, err := json.Marshal(billingmodel.GetInvoiceResponse{
		Invoice:     invoice,
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