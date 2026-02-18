package system

import (
	"context"
)

type Service interface {
	DB()
	Mux()
	RPC()
	Waiter()
	Logger()
}