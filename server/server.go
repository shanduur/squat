package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/shanduur/squat/server/api"
	"github.com/shanduur/squat/server/ui"
)

type WebServer struct {
	Server   *http.Server
	shutdown chan bool
}

func New(addr string) *WebServer {
	r := mux.NewRouter()

	rApi := r.PathPrefix("/api/v1").Subrouter()

	api.RegisterEndpoints(rApi)
	ui.RegisterEndpoints(r)

	return &WebServer{
		Server: &http.Server{
			Addr:    addr,
			Handler: r,
		},
		shutdown: make(chan bool),
	}
}

func (srv *WebServer) run() error {
	return srv.Server.ListenAndServe()
}

func (srv *WebServer) teardown(ctx context.Context) error {
	if err := srv.Server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (srv *WebServer) Run() error {
	term := make(chan os.Signal)
	fail := make(chan error)

	go func() {
		signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-term: // checks for termination signal
			// context with 30s timeout
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// all teardown process must complete within 30 seconds
			fail <- srv.teardown(ctx)

			return

		case <-srv.shutdown:
			// context with 30s timeout
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// all teardown process must complete within 30 seconds
			fail <- srv.teardown(ctx)

			return
		}
	}()

	// listenAndServe blocks our code from exit, but will produce ErrServerClosed when stopped
	if err := srv.run(); err != nil && err != http.ErrServerClosed {
		return err
	}

	// after server gracefully stopped, code proceeds here and waits for any error produced by teardown() process @ line 26
	return <-fail
}

func (srv *WebServer) Shutdown() {
	srv.shutdown <- true
}
