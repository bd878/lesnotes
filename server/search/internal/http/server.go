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
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	usersgateway "github.com/bd878/gallery/server/internal/gateway/users"
	sessionsgateway "github.com/bd878/gallery/server/internal/gateway/sessions"
	httphandler "github.com/bd878/gallery/server/search/internal/handler/http"
	controller "github.com/bd878/gallery/server/search/internal/controller/service"
)

type Config struct {
	Addr                string
	RpcAddr             string
	UsersServiceAddr    string
	SessionsServiceAddr string
}

type Server struct {
	*http.Server
	conf  Config
}

func New(conf Config) *Server {
	mux := http.NewServeMux()

	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	server := &Server{
		Server: &http.Server{
			Addr:    conf.Addr,
			Handler: mux,
		},
		conf: conf,
	}

	usersGateway := usersgateway.New(conf.UsersServiceAddr)
	sessionsGateway := sessionsgateway.New(conf.SessionsServiceAddr)

	ctrl := controller.New(controller.Config{RpcAddr: conf.RpcAddr})

	handler := httphandler.New(ctrl)

	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(logger.Default(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	mux.Handle("/search/v1/messages", middleware.Build(handler.SearchMessages))
	mux.Handle("/search/v1/files", middleware.Build(handler.SearchFiles))

	middleware.NoAuth()
	mux.Handle("/search/v1/status", middleware.Build(handler.GetStatus))

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(logger.Default(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	mux.Handle("/search/v2/messages", middleware.Build(handler.SearchMessagesJsonAPI))
	mux.Handle("/search/v2/files", middleware.Build(handler.SearchFilesJsonAPI))

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
