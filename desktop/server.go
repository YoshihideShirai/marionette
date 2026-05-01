package desktop

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

type localServer struct {
	URL    string
	server *http.Server
	errc   chan error
}

func startLocalServer(handler http.Handler) (*localServer, error) {
	if handler == nil {
		return nil, errors.New("desktop: handler is nil")
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	srv := &http.Server{Handler: handler}
	local := &localServer{
		URL:    fmt.Sprintf("http://%s/", listener.Addr().String()),
		server: srv,
		errc:   make(chan error, 1),
	}
	go func() {
		err := srv.Serve(listener)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		local.errc <- err
	}()
	return local, nil
}

func (s *localServer) Shutdown(ctx context.Context) error {
	if s == nil || s.server == nil {
		return nil
	}
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
	}
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return <-s.errc
}
