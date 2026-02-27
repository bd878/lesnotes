package loadbalance

import "github.com/bd878/gallery/server/internal/balancer"

func init() {
	balancer.RegisterResolver(Name)
	balancer.RegisterPicker(
		Name,
		[]string{"SaveMessage", "DeleteMessage", "PublishMessage", "PrivateMessage", "UpdateMessage"},
		[]string{"SearchMessages"},
	)
}

const Name = "search"
