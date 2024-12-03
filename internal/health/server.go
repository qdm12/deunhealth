package health

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	address string
	infoer  Logger
	handler http.Handler
}

func NewServer(address string, infoer Logger, healthcheck func() error) *Server {
	handler := newHandler(healthcheck)
	return &Server{
		address: address,
		infoer:  infoer,
		handler: handler,
	}
}

func (s *Server) Run(ctx context.Context) error {
	const readTimeout = time.Second
	server := http.Server{Addr: s.address, Handler: s.handler, ReadTimeout: readTimeout}
	shutdownErrCh := make(chan error)
	go func() {
		<-ctx.Done()
		const shutdownGraceDuration = 2 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGraceDuration)
		defer cancel()
		shutdownErrCh <- server.Shutdown(shutdownCtx) //nolint:contextcheck
	}()

	s.infoer.Info("listening on " + s.address)
	err := server.ListenAndServe()
	if err != nil && !errors.Is(ctx.Err(), context.Canceled) {
		return fmt.Errorf("health server crashed: %w", err)
	}

	if err := <-shutdownErrCh; err != nil {
		return fmt.Errorf("health server failed shutting down: %w", err)
	}

	return nil
}
