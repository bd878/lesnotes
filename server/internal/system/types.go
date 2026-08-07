package system

import (
	"database/sql"
	"net"
	"net/http"

	"github.com/soheilhy/cmux"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/waiter"
)

type Service interface {
	DB() *sql.DB
	Config() Config
	Pool() *pgxpool.Pool
	Mux() cmux.CMux
	ServeMux() *http.ServeMux
	HTTP() *http.Server
	JS() nats.JetStreamContext
	RPC() *grpc.Server
	Waiter() waiter.Waiter
	RaftListener() net.Listener
}