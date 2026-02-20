package http

import (
	"os"
	"fmt"
	"time"
	"context"
	"net/http"
	"golang.org/x/sync/errgroup"

	"github.com/bd878/gallery/server/internal/waiter"
	"github.com/bd878/gallery/server/internal/logger"
	usermodel "github.com/bd878/gallery/server/users/pkg/model"
	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	httphandler "github.com/bd878/gallery/server/messages/internal/handler/http"
	usersgateway "github.com/bd878/gallery/server/internal/gateway/users"
	sessionsgateway "github.com/bd878/gallery/server/internal/gateway/sessions"
	filesgateway "github.com/bd878/gallery/server/messages/internal/gateway/files/grpc"
	threadsgateway "github.com/bd878/gallery/server/messages/internal/gateway/threads/grpc"
	controller "github.com/bd878/gallery/server/messages/internal/controller/service"
)

type Config struct {
	Addr                string
	RpcAddr             string
	UsersServiceAddr    string
	FilesServiceAddr    string
	SessionsServiceAddr string
	ThreadsServiceAddr  string
}

type Server struct {
	*http.Server
	config Config
}

func New(cfg Config) *Server {
	mux := http.NewServeMux()

	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	usersGateway := usersgateway.New(cfg.UsersServiceAddr)
	filesGateway := filesgateway.New(cfg.FilesServiceAddr)
	sessionsGateway := sessionsgateway.New(cfg.SessionsServiceAddr)
	threadsGateway := threadsgateway.New(cfg.ThreadsServiceAddr)
	middleware = middleware.WithAuth(httpmiddleware.AuthBuilder(logger.Default(), usersGateway, sessionsGateway, usermodel.PublicUserID))

	messagesController := controller.NewMessagesController(controller.MessagesConfig{
		RpcAddr: cfg.RpcAddr,
	}, threadsGateway)
	translationsController := controller.NewTranslationsController(controller.TranslationsConfig{
		RpcAddr: cfg.RpcAddr,
	})
	handler := httphandler.New(messagesController, translationsController, filesGateway)

	mux.Handle("/messages/v1/send",        middleware.Build(handler.SendMessage))
	mux.Handle("/messages/v1/read_path",   middleware.Build(handler.ReadPath))
	mux.Handle("/messages/v1/read",        middleware.Build(handler.ReadMessages))
	mux.Handle("/messages/v1/update",      middleware.Build(handler.UpdateMessage))
	mux.Handle("/messages/v1/publish",     middleware.Build(handler.PublishMessages))
	mux.Handle("/messages/v1/private",     middleware.Build(handler.PrivateMessages))
	mux.Handle("/messages/v1/delete",      middleware.Build(handler.DeleteMessages))

	middleware.NoAuth()
	mux.Handle("/messages/v1/status",      middleware.Build(handler.GetStatus))

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(logger.Default(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	mux.Handle("/messages/v2/send",        middleware.Build(handler.SendMessageJsonAPI))
	mux.Handle("/messages/v2/read",        middleware.Build(handler.ReadMessagesJsonAPI))
	// TODO: /threads/v2/read
	mux.Handle("/messages/v2/read_path",   middleware.Build(handler.ReadPathJsonAPI))
	mux.Handle("/messages/v2/publish",     middleware.Build(handler.PublishMessagesJsonAPI))
	mux.Handle("/messages/v2/private",     middleware.Build(handler.PrivateMessagesJsonAPI))
	mux.Handle("/messages/v2/delete",      middleware.Build(handler.DeleteMessagesJsonAPI))
	mux.Handle("/messages/v2/update",      middleware.Build(handler.UpdateMessageJsonAPI))

	mux.Handle("/translations/v2/send",    middleware.Build(handler.SendTranslationJsonAPI))
	mux.Handle("/translations/v2/update",  middleware.Build(handler.UpdateTranslationJsonAPI))
	mux.Handle("/translations/v2/delete",  middleware.Build(handler.DeleteTranslationJsonAPI))
	mux.Handle("/translations/v2/read",    middleware.Build(handler.ReadTranslationJsonAPI))
	mux.Handle("/translations/v2/list",    middleware.Build(handler.ListTranslationsJsonAPI))

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
