package loadbalance

import "github.com/bd878/gallery/server/internal/balancer"

func init() {
	balancer.RegisterResolver(Name)
	balancer.RegisterPicker(
		Name,
		[]string{"Apply"},
		[]string{"Resolve", "Read", "Count", "CountMessages",
			"List", "ListMessages", "ReadParent"},
	)
}

const Name = "threads"
