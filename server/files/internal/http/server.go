package http

import (
	"os"
	"fmt"
	"time"
	"context"
	"net/http"
	"golang.org/x/sync/errgroup"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bd878/gallery/server/api/users"
	"github.com/bd878/gallery/server/api/sessions"
	"github.com/bd878/gallery/server/internal/waiter"
	"github.com/bd878/gallery/server/internal/di"
	"github.com/bd878/gallery/server/internal/rpc"
	sessionsloadbalance "github.com/bd878/gallery/server/db/sessions/pkg/loadbalance"
	usersloadbalance "github.com/bd878/gallery/server/db/users/pkg/loadbalance"
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
	container := di.New()
	mux := http.NewServeMux()

	container.AddSingleton("sessionsConn", func(c di.Container) (any, error) {
		return rpc.NewClient(
			fmt.Sprintf(
				"%s:///%s",
				sessionsloadbalance.Name,
				cfg.SessionsServiceAddr,
			),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	})
	container.AddSingleton("usersConn", func(c di.Container) (any, error) {
		return rpc.NewClient(
			fmt.Sprintf(
				"%s:///%s",
				usersloadbalance.Name,
				cfg.UsersServiceAddr,
			),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	})
	container.AddSingleton("usersClient", func(c di.Container) (any, error) {
		client := users.NewUsersClient(c.Get("usersConn").(*grpc.ClientConn))
		return client, nil
	})
	container.AddSingleton("sessionsClient", func(c di.Container) (any, error) {
		client := sessions.NewSessionsClient(c.Get("sessionsConn").(*grpc.ClientConn))
		return client, nil
	})

	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	usersGateway := usersgateway.New(container)
	sessionsGateway := sessionsgateway.New(container)
	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(usersGateway, sessionsGateway, usermodel.PublicUserID))

	grpcCtrl := controller.New(controller.Config{
		RpcAddr: cfg.RpcAddr,
	})
	handler := httphandler.New(grpcCtrl)

	mux.Handle("GET /files/v1/download",         middleware.Build(handler.DownloadFile))
	mux.Handle("GET /files/v1/read/{id}",        middleware.Build(handler.ReadFile))
	mux.Handle("POST /files/v1/upload",          middleware.Build(handler.UploadFile))

	middleware.NoAuth()
	mux.Handle("GET /liveness",           middleware.Build(handler.GetStatus))

	middleware.NoAuth().WithAuth(httpmiddleware.TokenAuthBuilder(usersGateway, sessionsGateway, usermodel.PublicUserID))
	mux.Handle("POST /files/v2/read_meta",       middleware.Build(handler.ReadFileMetaJsonAPI))
	mux.Handle("POST /files/v2/upload",          middleware.Build(handler.UploadFileJsonAPI))
	mux.Handle("DELETE /files/v2/delete",        middleware.Build(handler.DeleteFileJsonAPI))
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
