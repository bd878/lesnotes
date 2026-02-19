package system

import (
	"database/sql"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/waiter"
	"github.com/bd878/gallery/server/logger"
)

type Service interface {
	DB() *sql.DB
	Config() Config
	Pool() *pgxpool.Pool
	Mux() cmux.CMux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	Logger() logger.Logger
}