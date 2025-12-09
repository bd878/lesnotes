package http

import (
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"

	middleware "github.com/bd878/gallery/server/internal/middleware/http"
	usersmodel "github.com/bd878/gallery/server/users/pkg/model"
	billingmodel "github.com/bd878/gallery/server/billing/pkg/model"
	server "github.com/bd878/gallery/server/pkg/model"
)

func (h *Handler) GetPayment(w http.ResponseWriter, req *http.Request) (err error) {
	var paymentID int64

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
		id, err := strconv.Atoi(values.Get("id"))
		if err != nil {
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

		paymentID = int64(id)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error: &server.ErrorCode{
				Code:    server.CodeNoID,
				Explain: "no payment id",
			},
		})

		return err
	}

	payment, err := h.controller.GetPayment(req.Context(), paymentID, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(server.ServerResponse{
			Status: "error",
			Error:  &server.ErrorCode{
				Code:    server.CodeWrongFormat,
				Explain: "failed to get payment",
			},
		})

		return err
	}

	response, err := json.Marshal(billingmodel.GetPaymentResponse{
		Payment: payment,
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