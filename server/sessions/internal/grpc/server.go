package grpc

import (
	"net"
	"sync"
	"context"
	"google.golang.org/grpc"
	"github.com/soheilhy/cmux"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/logger"

	grpcmiddleware "github.com/bd878/gallery/server/internal/middleware/grpc"
	controller "github.com/bd878/gallery/server/sessions/internal/controller/sessions"
	repository "github.com/bd878/gallery/server/sessions/internal/repository/sqlite"
	grpchandler "github.com/bd878/gallery/server/sessions/internal/handler/grpc"
)

type Config struct {
	Addr            string
	DBPath          string
	NodeName        string
}

type Server struct {
	*grpc.Server
	config           Config
	mux              cmux.CMux
	listener         net.Listener
	grpcListener     net.Listener
}

func New(cfg Config) *Server {
	listener, err := net.Listen("tcp4", cfg.Addr)
	if err != nil {
		panic(err)
	}

	mux := cmux.New(listener)

	server := &Server{
		config:        cfg,
		mux:           mux,
		listener:      listener,
	}

	server.setupGRPC(logger.Default())

	return server
}

func (s *Server) setupGRPC(log *logger.Logger) {
	repo := repository.New(s.config.DBPath)
	control := controller.New(repo)
	handler := grpchandler.New(control)

	s.Server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.UnaryServerInterceptor(grpcmiddleware.LogBuilder()),
		),
	)

	api.RegisterSessionsServer(s.Server, handler)

	s.grpcListener = s.mux.Match(cmux.Any())
}

func (s *Server) Run(ctx context.Context, wg *sync.WaitGroup) {
	go s.Serve(s.grpcListener)
	defer s.mux.Close()
	s.mux.Serve()
	wg.Done()
}

func (s *Server) Addr() string {
	return s.listener.Addr().String()
}