package model

type (
	Invoice struct {
		ID          string        `json:"id"`
		UserID      int64         `json:"user_id"`
		Status      string        `json:"status"`
		Currency    string        `json:"currency"`
		Total       int64         `json:"total"`
		CreatedAt   string        `json:"created_at"`
		UpdatedAt   string        `json:"updated_at"`
	}

	Payment struct {
		ID          int64         `json:"id"`
		UserID      int64         `json:"user_id"`
		InvoiceID   string        `json:"invoice_id"`
		Status      string        `json:"status"`
		Currency    string        `json:"currency"`
		Total       int64         `json:"total"`
		CreatedAt   string        `json:"created_at"`
		UpdatedAt   string        `json:"updated_at"`
	}

	CreateInvoiceRequest struct {
		ID          string        `json:"string"`
		Currency    string        `json:"currency"`
		Total       int64         `json:"total"`
	}

	CreateInvoiceResponse struct {
		Description string        `json:"description"`
	}

	StartPaymentRequest struct {
		ID          int64         `json:"id"`
		InvoiceID   string        `json:"invoice_id"`
		Total       int64         `json:"total"`
		Currency    string        `json:"currency"`
	}

	StartPaymentResponse struct {
		Description string        `json:"description"`
	}

	GetInvoiceRequest struct {
		ID          string        `json:"id"`
	}

	GetInvoiceResponse struct {
		Invoice     Invoice       `json:"invoice"`
		Description string        `json:"description"`
	}

	GetPaymentRequest struct {
		Payment     Payment       `json:"payment"`
		Description string        `json:"description"`
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
		ID          string        `json:"id"`
	}

	RefundPaymentResponse struct {
		Description string        `json:"description"`
	}
)