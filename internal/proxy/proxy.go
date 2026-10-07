// Package proxy forwards local tcp connections to a database.
package proxy

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"
)

const dialTimeout = 10 * time.Second

// Serve accepts connections on l and pipes each one to target until ctx is canceled,
// then returns ctx's error once every connection has closed.
func Serve(ctx context.Context, l net.Listener, target string, logger *slog.Logger) error {
	stop := context.AfterFunc(ctx, func() { _ = l.Close() })
	defer stop()

	var wg sync.WaitGroup
	defer wg.Wait()

	for {
		conn, err := l.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return fmt.Errorf("stop serving: %w", ctx.Err())
			}

			return fmt.Errorf("accept connection: %w", err)
		}

		wg.Go(func() { pipe(ctx, conn, target, logger) })
	}
}

func pipe(ctx context.Context, client net.Conn, target string, logger *slog.Logger) {
	defer func() { _ = client.Close() }()

	d := net.Dialer{Timeout: dialTimeout}

	upstream, err := d.DialContext(ctx, "tcp", target)
	if err != nil {
		logger.Warn("dial database", "target", target, "error", err)

		return
	}

	defer func() { _ = upstream.Close() }()

	stop := context.AfterFunc(ctx, func() {
		_ = client.Close()
		_ = upstream.Close()
	})
	defer stop()

	var wg sync.WaitGroup

	wg.Go(func() { copyHalf(upstream, client) })
	copyHalf(client, upstream)
	wg.Wait()
}

// copyHalf closes only the write side when done so the other direction keeps flowing.
func copyHalf(dst, src net.Conn) {
	_, _ = io.Copy(dst, src)

	if tcp, ok := dst.(*net.TCPConn); ok {
		_ = tcp.CloseWrite()

		return
	}

	_ = dst.Close()
}
