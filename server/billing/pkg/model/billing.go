package model

import "encoding/json"

type (
	Invoice struct {
		ID          string            `json:"id"`
		UserID      int64             `json:"user_id"`
		Status      string            `json:"status"`
		Currency    string            `json:"currency"`
		Total       int64             `json:"total"`
		CreatedAt   string            `json:"created_at"`
		UpdatedAt   string            `json:"updated_at"`
		Metadata    json.RawMessage   `json:"metadata"`
	}

	Payment struct {
		ID          int64             `json:"id"`
		UserID      int64             `json:"user_id"`
		InvoiceID   string            `json:"invoice_id"`
		Status      string            `json:"status"`
		Currency    string            `json:"currency"`
		Total       int64             `json:"total"`
		CreatedAt   string            `json:"created_at"`
		UpdatedAt   string            `json:"updated_at"`
		Metadata    json.RawMessage   `json:"metadata"`
	}

	CreateInvoiceRequest struct {
		ID          string            `json:"string"`
		Currency    string            `json:"currency"`
		Total       int64             `json:"total"`
		Metadata    json.RawMessage   `json:"metadata"`
	}

	CreateInvoiceResponse struct {
		Description string            `json:"description"`
	}

	StartPaymentRequest struct {
		ID          int64             `json:"id"`
		InvoiceID   string            `json:"invoice_id"`
		Total       int64             `json:"total"`
		Currency    string            `json:"currency"`
		Metadata    json.RawMessage   `json:"metadata"`
	}

	StartPaymentResponse struct {
		Description string        `json:"description"`
	}

	GetInvoiceRequest struct {
		ID          string        `json:"id"`
	}

	GetInvoiceResponse struct {
		Invoice     *Invoice      `json:"invoice"`
	}

	GetPaymentRequest struct {
		ID          int64         `json:"id"`
	}

	GetPaymentResponse struct {
		Payment     *Payment      `json:"payment"`
	}

	CancelPaymentRequest struct {
		ID          int64         `json:"id"`
	}

	CancelPaymentResponse struct {
		Description string        `json:"description"`
	}

	FulfillPaymentRequest struct {
		ID          string        `json:"id"`
	}

	FulfillPaymentResponse struct {
		Description string        `json:"description"`
	}

	RefundPaymentRequest struct {
		ID          int64         `json:"id"`
	}

	RefundPaymentResponse struct {
		Description string        `json:"description"`
	}
)