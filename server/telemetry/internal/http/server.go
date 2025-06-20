package http

import (
	"context"
	"sync"
	"net/http"

	httpmiddleware "github.com/bd878/gallery/server/internal/middleware/http"
	httphandler "github.com/bd878/gallery/server/telemetry/internal/handler/http"
)

type Config struct {
	Addr             string
}

type Server struct {
	*http.Server
	config Config
}

func New(cfg Config) *Server {
	mux := http.NewServeMux()

	middleware := httpmiddleware.NewBuilder().WithLog(httpmiddleware.Log)

	handler := httphandler.New()

	mux.Handle("/telemetry/v1/send", middleware.Build(handler.SendLog))
	mux.Handle("/telemetry/v1/status", middleware.Build(handler.GetStatus))

	server := &Server{
		Server: &http.Server{
			Addr: cfg.Addr,
			Handler: mux,
		},
		config: cfg,
	}

	return server
}

func (s *Server) ListenAndServe(_ context.Context, wg *sync.WaitGroup) {
	err := s.Server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	wg.Done()
}