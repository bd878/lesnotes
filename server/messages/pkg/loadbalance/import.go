package loadbalance

import "github.com/bd878/gallery/server/internal/balancer"

func init() {
	balancer.RegisterResolver(Name)
	balancer.RegisterPicker(
		Name,
		[]string{"SaveMessage", "DeleteMessages", "DeleteUserMessages",
			"PublishMessages", "PrivateMessages", "UpdateMessage",
			"SaveTranslation", "UpdateTranslation", "DeleteTranslation"},
		[]string{"ReadMessages", "ReadMessage", "ReadTranslation", "ListTranslations",
			"ReadBatchMessages"},
	)
}

const Name = "messages"
