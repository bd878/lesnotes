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
	httphandler "github.com/bd878/gallery/server/billing/internal/handler/http"
	controller "github.com/bd878/gallery/server/billing/internal/controller/service"
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

	middleware.NoAuth()
	mux.Handle("GET   /billing/v1/status", middleware.Build(handler.GetStatus))

	middleware.WithAuth(httpmiddleware.TokenAuthBuilder(logger.Default(), usersGateway, sessionsGateway, usermodel.PublicUserID))
	mux.Handle("POST  /billing/v2/invoices", middleware.Build(handler.CreateInvoiceJsonAPI))
	mux.Handle("GET   /billing/v2/invoices", middleware.Build(handler.GetInvoiceJsonAPI))
	mux.Handle("POST  /billing/v2/payments", middleware.Build(handler.StartPaymentJsonAPI))
	mux.Handle("GET   /billing/v2/payments", middleware.Build(handler.GetPaymentJsonAPI))
	mux.Handle("POST  /billing/v2/cancel",   middleware.Build(handler.CancelPaymentJsonAPI))
	mux.Handle("POST  /billing/v2/refund",   middleware.Build(handler.RefundPaymentJsonAPI))

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
