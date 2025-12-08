package model

import (
	"github.com/bd878/gallery/server/api"
)

func InvoiceFromProto(proto *api.Invoice) *Invoice {
	return &Invoice{
		ID:             proto.Id,
		UserID:         proto.UserId,
		Status:         proto.Status,
		Currency:       proto.Currency,
		Total:          proto.Total,
		CreatedAt:      proto.CreatedAt,
		UpdatedAt:      proto.UpdatedAt,
		Metadata:       proto.Metadata,
	}
}

func InvoiceToProto(invoice *Invoice) *api.Invoice {
	return &api.Invoice{
		Id:             invoice.ID,
		UserId:         invoice.UserID,
		Currency:       invoice.Currency,
		Total:          invoice.Total,
		Status:         invoice.Status,
		CreatedAt:      invoice.CreatedAt,
		UpdatedAt:      invoice.UpdatedAt,
		Metadata:       invoice.Metadata,
	}
}

func PaymentFromProto(proto *api.Payment) *Payment {
	return &Payment{
		ID:             proto.Id,
		UserID:         proto.UserId,
		InvoiceID:      proto.InvoiceId,
		Status:         proto.Status,
		Currency:       proto.Currency,
		Total:          proto.Total,
		CreatedAt:      proto.CreatedAt,
		UpdatedAt:      proto.UpdatedAt,
		Metadata:       proto.Metadata,
	}
}

func PaymentToProto(payment *Payment) *api.Payment {
	return &api.Payment{
		Id:             payment.ID,
		UserId:         payment.UserID,
		InvoiceId:      payment.InvoiceID,
		Currency:       payment.Currency,
		Total:          payment.Total,
		Status:         payment.Status,
		CreatedAt:      payment.CreatedAt,
		UpdatedAt:      payment.UpdatedAt,
		Metadata:       payment.Metadata,
	}
}

func MapInvoicesToProto(mapper (func(*Invoice) *api.Invoice), invoices []*Invoice) []*api.Invoice {
	res := make([]*api.Invoice, len(invoices))
	for i, invoice := range invoices {
		res[i] = mapper(invoice)
	}
	return res
}

func MapInvoicesFromProto(mapper (func(*api.Invoice) *Invoice), invoices []*api.Invoice) []*Invoice {
	res := make([]*Invoice, len(invoices))
	for i, invoice := range invoices {
		res[i] = mapper(invoice)
	}
	return res
}

func MapPaymentsToProto(mapper (func(*Payment) *api.Payment), payments []*Payment) []*api.Payment {
	res := make([]*api.Payment, len(payments))
	for i, payment := range payments {
		res[i] = mapper(payment)
	}
	return res
}

func MapPaymentsFromProto(mapper (func(*api.Payment) *Payment), payments []*api.Payment) []*Payment {
	res := make([]*Payment, len(payments))
	for i, payment := range payments {
		res[i] = mapper(payment)
	}
	return res
}