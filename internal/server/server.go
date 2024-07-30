package server

import (
	"context"
	"net/http"
	"time"

	"github.com/VladislavLisovenko/antibruteforce/internal/app"
	"github.com/VladislavLisovenko/antibruteforce/internal/logger"
)

type server struct {
	addr   string
	srv    *http.Server
	app    app.Application
	logger logger.Logger
}

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

func New(addr string, a app.Application, l logger.Logger) Server {
	return &server{
		addr:   addr,
		app:    a,
		logger: l,
	}
}

func (s *server) Start(ctx context.Context) error {
	s.logger.Info("Start server on " + s.addr)
	s.srv = &http.Server{
		Addr:              s.addr,
		Handler:           handler(ctx, s.app, s.logger),
		ReadHeaderTimeout: 10 * time.Second,
	}

	err := s.srv.ListenAndServe()

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (s *server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func handler(ctx context.Context, app app.Application, logg logger.Logger) http.Handler {
	service := NewHandlers(app, logg)

	h := loggingMiddleware(ctx, service.Handlers(ctx), logg)

	return h
}
