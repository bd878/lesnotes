package system

import (
	"fmt"
	"io/fs"
	"os"
	"net"
	"time"
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/nats-io/nats.go"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"golang.org/x/sync/errgroup"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/internal/logger"
	"github.com/bd878/gallery/server/internal/waiter"
	grpcmiddleware "github.com/bd878/gallery/server/internal/middleware/grpc"
)

type Config struct {
	RpcAddr           string
	NodeName          string
	LogLevel          string
	SkipCaller        int
	NatsAddr          string
	PGConn            string
	GooseTableName    string
}

type System struct {
	cfg          Config
	pool         *pgxpool.Pool
	db           *sql.DB
	nc           *nats.Conn
	rpc          *grpc.Server
	logger       *logger.Logger
	listener     net.Listener
	waiter       waiter.Waiter
	mux          cmux.CMux
}

func NewSystem(cfg Config) (*System, error) {
	s := &System{cfg: cfg}

	if err := s.initDB(); err != nil {
		return nil, err
	}

	if err := s.initPool(); err != nil {
		return nil, err
	}

	if err := s.initNats(); err != nil {
		return nil, err
	}

	if err := s.initMux(); err != nil {
		return nil, err
	}

	s.initRPC()
	s.initLogger()
	s.initWaiter()

	return s, nil
}

func (s *System) Config() Config {
	return s.cfg
}

func(s *System) Listener() net.Listener {
	return s.listener
}

func (s *System) initDB() (err error) {
	s.db, err = sql.Open("pgx", s.cfg.PGConn)
	return
}

func (s *System) DB() *sql.DB {
	return s.db
}

func (s *System) initPool() (err error) {
	s.pool, err = pgxpool.New(context.Background(), s.cfg.PGConn)
	return
}

func (s *System) Pool() *pgxpool.Pool {
	return s.pool
}

func (s *System) initNats() (err error) {
	s.nc, err = nats.Connect(s.cfg.NatsAddr)
	return
}

func (s *System) initMux() (err error) {
	s.listener, err = net.Listen("tcp4", s.cfg.RpcAddr)
	if err != nil {
		return
	}

	s.mux = cmux.New(s.listener)

	return
}

func (s *System) Mux() cmux.CMux {
	return s.mux
}

func (s *System) MigrateDB(fs fs.FS) error {
	goose.SetVerbose(true)
	goose.SetTableName(s.cfg.GooseTableName)

	goose.SetBaseFS(fs)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(s.db, "."); err != nil {
		return err
	}

	return nil
}

func (s *System) ResetDB() error {
	return goose.Reset(s.db, ".")
}

func (s *System) initLogger() {
	s.logger = logger.New(logger.Config{
		NodeName:   s.cfg.NodeName,
		LogLevel:   s.cfg.LogLevel,
		SkipCaller: s.cfg.SkipCaller,
	})
}

func (s *System) Logger() *logger.Logger {
	return s.logger
}

func (s *System) initWaiter() {
	s.waiter = waiter.New(waiter.CatchSignals())
}

func (s *System) Waiter() waiter.Waiter {
	return s.waiter
}

func (s *System) initRPC() {
	s.rpc = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.UnaryServerInterceptor(grpcmiddleware.LogBuilder()),
		),
	)
}

func (s *System) RPC() *grpc.Server {
	return s.rpc
}

func (s *System) WaitForRPC(ctx context.Context) (err error) {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		fmt.Fprintf(os.Stdout, "rpc server started %s\n", s.Addr())
		defer fmt.Fprintln(os.Stdout, "rpc server shutdown")
		if err := s.RPC().Serve(s.Mux().Match(cmux.Any())); err != nil {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Fprintln(os.Stdout, "rpc server to be shutdown")
		stopped := make(chan struct{})
		go func() {
			s.RPC().GracefulStop()
			close(stopped)
		}()
		timeout := time.NewTimer(5*time.Second)
		select {
		case <-timeout.C:
			s.RPC().Stop()
			return fmt.Errorf("rpc server failed to stop gracefully")
		case <-stopped:
			return nil
		}
	})
	group.Go(func() error {
		fmt.Fprintln(os.Stdout, "mux serve")
		s.mux.Serve()
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Fprintln(os.Stdout, "mux close")
		s.mux.Close()
		return nil
	})

	return group.Wait()
}

func (s *System) WaitForPool(ctx context.Context) (err error) {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		<-gCtx.Done()
		fmt.Fprintln(os.Stdout, "closing pgpool connections")
		s.pool.Close()
		return nil
	})

	return group.Wait()
}

func (s *System) WaitForStream(ctx context.Context) error {
	closed := make(chan struct{})
	s.nc.SetClosedHandler(func (*nats.Conn) {
		close(closed)
	})
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Fprintln(os.Stdout, "message stream started")
		defer fmt.Fprintln(os.Stdout, "message stream stopped")
		<-closed
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		return s.nc.Drain()
	})
	return group.Wait()
}

func (s *System) Addr() string {
	return s.listener.Addr().String()
}