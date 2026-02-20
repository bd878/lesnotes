package system

import (
	"database/sql"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"github.com/nats-io/nats.go"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/waiter"
	"github.com/bd878/gallery/server/internal/logger"
)

type Service interface {
	DB() *sql.DB
	Config() Config
	Pool() *pgxpool.Pool
	Mux() cmux.CMux
	Nats() *nats.Conn
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() *logger.Logger
}