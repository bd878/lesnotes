package model

import (
	"encoding/json"
	"github.com/bd878/gallery/server/api/billing"
)

func InvoiceFromProto(proto *billing.Invoice) (*Invoice, error) {
	cart, err := CartFromProto(proto.Cart)
	if err != nil {
		return nil, err
	}

	return &Invoice{
		ID:             proto.Id,
		UserID:         proto.UserId,
		Status:         proto.Status,
		Total:          proto.Total,
		CreatedAt:      proto.CreatedAt,
		UpdatedAt:      proto.UpdatedAt,
		Metadata:       proto.Metadata,
		Cart:           cart,
	}, nil
}

func InvoiceToProto(invoice *Invoice) (*billing.Invoice, error) {
	cart, err := CartToProto(invoice.Cart)
	if err != nil {
		return nil, err
	}

	return &billing.Invoice{
		Id:             invoice.ID,
		UserId:         invoice.UserID,
		Total:          invoice.Total,
		Status:         invoice.Status,
		CreatedAt:      invoice.CreatedAt,
		UpdatedAt:      invoice.UpdatedAt,
		Metadata:       invoice.Metadata,
		Cart:           cart,
	}, nil
}

func PaymentFromProto(proto *billing.Payment) *Payment {
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

func PaymentToProto(payment *Payment) *billing.Payment {
	return &billing.Payment{
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

func CartFromProto(proto *billing.Cart) (*Cart, error) {
	items, err := MapCartItemsFromProto(CartItemFromProto, proto.Items)
	if err != nil {
		return nil, err
	}

	return &Cart{
		Items: items,
	}, nil
}

func CartToProto(cart *Cart) (*billing.Cart, error) {
	items, err := MapCartItemsToProto(CartItemToProto, cart.Items)
	if err != nil {
		return nil, err
	}

	return &billing.Cart{
		Items: items,
	}, nil
}

func CartItemFromProto(item *billing.CartItem) (*CartItem, error) {
	switch v := item.Item.(type) {
	case *billing.CartItem_Premium:
		premium, err := json.Marshal(&PremiumItem{
			ExpiresAt:   v.Premium.ExpiresAt,
			Cost:        v.Premium.Cost,
			Discount:    v.Premium.Discount,
			Currency:    v.Premium.Currency,
		})
		if err != nil {
			return nil, err
		}

		return &CartItem{
			Type: "premium",
			Item: json.RawMessage(premium),
		}, nil
	default:
		return &CartItem{Type: "unknown"}, nil
	}
}

func CartItemToProto(item *CartItem) (*billing.CartItem, error) {
	switch item.Type {
	case "premium":
		var premium PremiumItem
		err := json.Unmarshal(item.Item, &premium)
		if err != nil {
			return nil, err
		}

		return &billing.CartItem{
			Item: &billing.CartItem_Premium{
				Premium: &billing.Premium{
					ExpiresAt: premium.ExpiresAt,
					Cost:      premium.Cost,
					Discount:  premium.Discount,
					Currency:  premium.Currency,
				},
			},
		}, nil
	default:
		return &billing.CartItem{}, nil
	}
}

func MapInvoicesToProto(mapper (func(*Invoice) *billing.Invoice), invoices []*Invoice) []*billing.Invoice {
	res := make([]*billing.Invoice, len(invoices))
	for i, invoice := range invoices {
		res[i] = mapper(invoice)
	}
	return res
}

func MapInvoicesFromProto(mapper (func(*billing.Invoice) *Invoice), invoices []*billing.Invoice) []*Invoice {
	res := make([]*Invoice, len(invoices))
	for i, invoice := range invoices {
		res[i] = mapper(invoice)
	}
	return res
}

func MapPaymentsToProto(mapper (func(*Payment) *billing.Payment), payments []*Payment) []*billing.Payment {
	res := make([]*billing.Payment, len(payments))
	for i, payment := range payments {
		res[i] = mapper(payment)
	}
	return res
}

func MapPaymentsFromProto(mapper (func(*billing.Payment) *Payment), payments []*billing.Payment) []*Payment {
	res := make([]*Payment, len(payments))
	for i, payment := range payments {
		res[i] = mapper(payment)
	}
	return res
}

func MapCartItemsToProto(mapper (func(*CartItem) (*billing.CartItem, error)), items []*CartItem) (res []*billing.CartItem, err error) {
	res = make([]*billing.CartItem, len(items))
	for i, item := range items {
		res[i], err = mapper(item)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func MapCartItemsFromProto(mapper (func(*billing.CartItem) (*CartItem, error)), items []*billing.CartItem) (res []*CartItem, err error) {
	res = make([]*CartItem, len(items))
	for i, item := range items {
		res[i], err = mapper(item)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
