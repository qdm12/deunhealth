package health

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type ServerRunner interface {
	Run(ctx context.Context) error
}

type Server struct {
	address string
	infoer  Infoer
	handler http.Handler
}

type Infoer interface {
	Info(s string)
}

func NewServer(address string, infoer Infoer, healthcheck func() error) *Server {
	handler := newHandler(healthcheck)
	return &Server{
		address: address,
		infoer:  infoer,
		handler: handler,
	}
}

var (
	ErrCrashed  = errors.New("server crashed")
	ErrShutdown = errors.New("server could not be shutdown")
)

func (s *Server) Run(ctx context.Context) error {
	server := http.Server{Addr: s.address, Handler: s.handler}
	shutdownErrCh := make(chan error)
	go func() {
		<-ctx.Done()
		const shutdownGraceDuration = 2 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownGraceDuration)
		defer cancel()
		shutdownErrCh <- server.Shutdown(shutdownCtx)
	}()

	s.infoer.Info("listening on " + s.address)
	err := server.ListenAndServe()
	if err != nil && !errors.Is(ctx.Err(), context.Canceled) { // server crashed
		return fmt.Errorf("%w: %s", ErrCrashed, err)
	}

	if err := <-shutdownErrCh; err != nil {
		return fmt.Errorf("%w: %s", ErrShutdown, err)
	}

	return nil
}
