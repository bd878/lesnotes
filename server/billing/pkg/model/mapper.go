package model

import (
	"encoding/json"
	"github.com/bd878/gallery/server/api"
)

func InvoiceFromProto(proto *api.Invoice) (*Invoice, error) {
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

func InvoiceToProto(invoice *Invoice) (*api.Invoice, error) {
	cart, err := CartToProto(invoice.Cart)
	if err != nil {
		return nil, err
	}

	return &api.Invoice{
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

func CartFromProto(proto *api.Cart) (*Cart, error) {
	items, err := MapCartItemsFromProto(CartItemFromProto, proto.Items)
	if err != nil {
		return nil, err
	}

	return &Cart{
		Items: items,
	}, nil
}

func CartToProto(cart *Cart) (*api.Cart, error) {
	items, err := MapCartItemsToProto(CartItemToProto, cart.Items)
	if err != nil {
		return nil, err
	}

	return &api.Cart{
		Items: items,
	}, nil
}

func CartItemFromProto(item *api.CartItem) (*CartItem, error) {
	switch v := item.Item.(type) {
	case *api.CartItem_Premium:
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

func CartItemToProto(item *CartItem) (*api.CartItem, error) {
	switch item.Type {
	case "premium":
		var premium PremiumItem
		err := json.Unmarshal(item.Item, &premium)
		if err != nil {
			return nil, err
		}

		return &api.CartItem{
			Item: &api.CartItem_Premium{
				Premium: &api.Premium{
					ExpiresAt: premium.ExpiresAt,
					Cost:      premium.Cost,
					Discount:  premium.Discount,
					Currency:  premium.Currency,
				},
			},
		}, nil
	default:
		return &api.CartItem{}, nil
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

func MapCartItemsToProto(mapper (func(*CartItem) (*api.CartItem, error)), items []*CartItem) (res []*api.CartItem, err error) {
	res = make([]*api.CartItem, len(items))
	for i, item := range items {
		res[i], err = mapper(item)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func MapCartItemsFromProto(mapper (func(*api.CartItem) (*CartItem, error)), items []*api.CartItem) (res []*CartItem, err error) {
	res = make([]*CartItem, len(items))
	for i, item := range items {
		res[i], err = mapper(item)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
