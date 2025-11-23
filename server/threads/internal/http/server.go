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
	httphandler "github.com/bd878/gallery/server/threads/internal/handler/http"
	controller "github.com/bd878/gallery/server/threads/internal/controller/service"
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
	mux.Handle("/threads/v1/publish", middleware.Build(handler.PublishThread))
	mux.Handle("/threads/v1/private", middleware.Build(handler.PrivateThread))
	mux.Handle("/threads/v1/reorder", middleware.Build(handler.ReorderThread))

	middleware.NoAuth()
	mux.Handle("/threads/v1/status", middleware.Build(handler.GetStatus))

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(logger.Default(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	mux.Handle("/threads/v2/read", middleware.Build(handler.ReadThreadJsonAPI))
	mux.Handle("/threads/v2/list", middleware.Build(handler.ListThreadsJsonAPI))
	mux.Handle("/threads/v2/resolve", middleware.Build(handler.ResolveThreadJsonAPI))
	mux.Handle("/threads/v2/create", middleware.Build(handler.CreateThreadJsonAPI))
	mux.Handle("/threads/v2/delete", middleware.Build(handler.DeleteThreadJsonAPI))
	mux.Handle("/threads/v2/reorder", middleware.Build(handler.ReorderThreadJsonAPI))

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
