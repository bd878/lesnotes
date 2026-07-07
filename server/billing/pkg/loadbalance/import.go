package loadbalance

import "github.com/bd878/gallery/server/internal/balancer"

func init() {
	balancer.RegisterResolver(Name)
	balancer.RegisterPicker(
		Name,
		[]string{"CreateInvoice", "StartPayment", "ProceedPayment",
			"CancelPayment", "RefundPayment", "Apply"},
		[]string{"GetPayment", "GetInvoice"},
	)
}

const Name = "billing"
