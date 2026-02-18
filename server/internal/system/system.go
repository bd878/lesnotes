package system

import (
	"fmt"
	"fs"
	"os"
	"time"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/waiter"
)

type Config struct {
	NodeName    string
	LogLevel    string
	SkipCaller  int
	NatsAddr    string
	PGConn      string
}

type System struct {
	cfg       Config
	pool      *pgxpool.Pool
	db        *sql.DB // for migrations (TODO: get single connection from pool)
	nc        *nats.Conn
	logger    *logger.Logger
	waiter    waiter.Waiter
	mux       cmux.CMux
}

func NewSystem(cfg Config) (*System, error) {
	s := &System{cfg: cfg}

	if err := s.initDB(); err != nil {
		return nil, err
	}

	if err := s.initNats(); err != nil {
		return nil, err
	}

	s.initLogger()
	s.initWaiter()

	return s, nil
}

func (s *System) initDB() (err error) {
	s.pool, err = pgxpool.New(context.Background(), s.cfg.PGConn)
	if err != nil {
		return
	}
	s.db, err = sql.Open("pgx", s.cfg.PGConn)
	return
}

func (s *System) initNats() (err error) {
	s.nc, err = nats.Connect(s.cfg.NatsAddr)
	return
}

func (s *Service) MigrateDB(fs fs.FS) error {
	goose.SetVerbose(true)

	goose.SetBaseFS(fs)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(s.db, "."); err != nil {
		return err
	}

	return nil
}

func (s *Service) ResetDB() error {
	return goose.Reset(s.db, ".")
}

func (s *Serivce) DB() *sql.DB {
	return s.db
}

func (s *Service) Pool() *pgxpool.Pool {
	return s.pool
}

func (s *Service) initLogger() {
	s.logger = logger.New({
		NodeName:   s.cfg.NodeName,
		LogLevel:   s.cfg.LogLevel,
		SkipCaller: s.cfg.SkipCaller,
	})
}

func (s *Service) Logger() *logger.Logger {
	return s.logger
}

func (s *Service) initWaiter() {
	s.waiter = waiter.New(waiter.CatchSignals())
}

func (s *Service) Waiter() waiter.Waiter {
	return s.waiter
}

func (s *Service) WaitForRPC(ctx context.Context) (err error) {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		fmt.Fprintf(os.Stdout, "rpc server started %s\n", s.Addr())
		defer fmt.Fprintln(os.Stdout, "rpc server shutdown")
		if err := s.Serve(s.grpcListener); err != nil {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Fprintln(os.Stdout, "rpc server to be shutdown")
		stopped := make(chan struct{})
		go func() {
			s.GracefulStop()
			close(stopped)
		}()
		timeout := time.NewTimer(5*time.Second)
		select {
		case <-timeout.C:
			s.Stop()
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

func (s *Service) WaitForPool(ctx context.Context) (err error) {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		<-gCtx.Done()
		fmt.Fprintln(os.Stdout, "closing pgpool connections")
		s.pool.Close()
		return nil
	})

	return group.Wait()
}

func (s *Service) WaitForStream(ctx context.Context) error {
	closed := make(chan struct{})
	s.nc.SetClosedHandler(func (*nats.Conn) {
		close(closed)
	})
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Fprintln(os.Stdout, "messsage stream started")
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
