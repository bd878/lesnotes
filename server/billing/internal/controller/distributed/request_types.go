package distributed

type RequestType uint16

const (
	AppendInvoiceRequest RequestType = iota
	AppendPaymentRequest
	ProceedPaymentRequest
	PayInvoiceRequest
	CancelPaymentRequest
	CancelInvoiceRequest
	RefundPaymentRequest
)