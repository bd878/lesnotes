package loadbalance

import "github.com/bd878/gallery/server/internal/balancer"

func init() {
	balancer.RegisterResolver(Name)
	balancer.RegisterPicker(
		Name,
		[]string{"Create", "Delete", "Publish",
			"Private", "Update", "reorder"},
		[]string{"Resolve", "Read", "Count", "List"},
	)
}

const Name = "threads"
