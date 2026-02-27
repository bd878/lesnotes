package loadbalance

import "github.com/bd878/gallery/server/internal/balancer"

func init() {
	balancer.RegisterResolver(Name)
	balancer.RegisterPicker(
		Name,
		[]string{"Create", "Delete", "Publish",
			"Private", "Update", "Reorder"},
		[]string{"Resolve", "Read", "Count", "List"},
	)
}

const Name = "threads"
