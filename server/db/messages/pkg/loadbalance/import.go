package loadbalance

import "github.com/bd878/gallery/server/internal/balancer"

func init() {
	balancer.RegisterResolver(Name)
	balancer.RegisterPicker(
		Name,
		[]string{"Apply"},
		[]string{"ReadMessages", "ReadMessage", "ReadTranslation", "ListTranslations",
			"ReadBatchMessages", "ReadComment", "ListComments"},
	)
}

const Name = "messages"
