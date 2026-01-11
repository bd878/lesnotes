package http

import (
	"os"
	"fmt"
	"time"
	"context"
	"net/http"
	"golang.org/x/sync/errgroup"

	"github.com/bd878/gallery/server/waiter"
	"github.com/bd878/gallery/server/logger"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	usersgateway "github.com/bd878/gallery/server/internal/gateway/users"
	sessionsgateway "github.com/bd878/gallery/server/internal/gateway/sessions"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	httphandler "github.com/bd878/gallery/server/files/internal/handler/http"
	controller "github.com/bd878/gallery/server/files/internal/controller/service"
)

type Config struct {
	Addr                string
	RpcAddr             string
	UsersServiceAddr    string
	SessionsServiceAddr string
	DataPath            string
}

type Server struct {
	*http.Server
	config Config
}

func New(cfg Config) *Server {
	mux := http.NewServeMux()

	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	usersGateway := usersgateway.New(cfg.UsersServiceAddr)
	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)
	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(logger.Default(), usersGateway, sessionsGateway, usermodel.PublicUserID))

	grpcCtrl := controller.New(controller.Config{
		RpcAddr: cfg.RpcAddr,
	})
	handler := httphandler.New(grpcCtrl)

	mux.Handle("GET /files/v1/download",         middleware.Build(handler.DownloadFile))
	mux.Handle("GET /files/v1/read/{name}",      middleware.Build(handler.ReadFile))
	mux.Handle("POST /files/v1/upload",          middleware.Build(handler.UploadFile))

	middleware.NoAuth()
	mux.Handle("GET /files/v1/status",           middleware.Build(handler.GetStatus))

	middleware.NoAuth().WithAuth(httpmiddleware.TokenAuthBuilder(logger.Default(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	mux.Handle("POST /files/v2/upload",          middleware.Build(handler.UploadFileV2))
	mux.Handle("POST /files/v2/list",            middleware.Build(handler.ListFilesJsonAPI))
	mux.Handle("POST /files/v2/publish",         middleware.Build(handler.PublishFileJsonAPI))
	mux.Handle("POST /files/v2/private",         middleware.Build(handler.PrivateFileJsonAPI))

	server := &Server{
		Server: &http.Server{
			Addr: cfg.Addr,
			Handler: mux,
		},
		config: cfg,
	}

	return server
}

func (s *Server) Run(ctx context.Context) (err error) {
	waiter := waiter.New(waiter.CatchSignals())

	waiter.Add(s.WaitForServer)

	return waiter.Wait()
}

func (s *Server) WaitForServer(ctx context.Context) (err error) {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		fmt.Fprintf(os.Stdout, "http server started %s\n", s.Addr)
		defer fmt.Fprintln(os.Stdout, "http server shutdown")
		if err := s.ListenAndServe(); err != nil {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Fprintln(os.Stdout, "http server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "http server failed to stop gracefully")
			return err
		}
		return nil
	})

	return group.Wait()
}
