package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/aelpxy/dbctl/internal/api"
	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/spf13/cobra"
)

const (
	readHeaderTimeout = 10 * time.Second
	shutdownTimeout   = 5 * time.Second
)

func (a *app) newHTTPCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "http <host:port>",
		Short:   "Start the API and serve it",
		Example: "  dbctl http localhost:5000",
		Aliases: []string{"serve"},
		Args:    cobra.ExactArgs(1),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			return serve(ctx, c, args[0])
		}),
	}
}

func serve(ctx context.Context, c *docker.Client, addr string) error {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	lc := net.ListenConfig{}

	l, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	server := &http.Server{
		Handler:           api.NewHandler(c, logger),
		ReadHeaderTimeout: readHeaderTimeout,
		BaseContext:       func(net.Listener) context.Context { return ctx },
	}

	errs := make(chan error, 1)

	go func() { errs <- server.Serve(l) }()

	logger.Info("server started", "address", "http://"+l.Addr().String())

	select {
	case err := <-errs:
		return fmt.Errorf("serve http: %w", err)
	case <-ctx.Done():
	}

	logger.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shut down server: %w", err)
	}

	if err := <-errs; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve http: %w", err)
	}

	logger.Info("server stopped")

	return nil
}
