package proxy_test

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"testing"

	"github.com/aelpxy/dbctl/internal/proxy"
)

func listen(t *testing.T) net.Listener {
	t.Helper()

	lc := net.ListenConfig{}

	l, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	return l
}

// echo serves a line-echoing upstream standing in for the database.
func echo(t *testing.T) string {
	t.Helper()

	l := listen(t)
	t.Cleanup(func() { _ = l.Close() })

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}

			go func() {
				defer func() { _ = conn.Close() }()

				_, _ = io.Copy(conn, conn)
			}()
		}
	}()

	return l.Addr().String()
}

func TestServe(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	upstream := echo(t)
	l := listen(t)
	done := make(chan error, 1)

	go func() { done <- proxy.Serve(ctx, l, upstream, slog.New(slog.DiscardHandler)) }()

	d := net.Dialer{}

	conn, err := d.DialContext(t.Context(), "tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("dial proxy: %v", err)
	}

	if _, err := io.WriteString(conn, "ping\n"); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	if got != "ping\n" {
		t.Errorf("echo through proxy = %q, want %q", got, "ping\n")
	}

	cancel()

	if err := <-done; !errors.Is(err, context.Canceled) {
		t.Errorf("Serve() error = %v, want context.Canceled", err)
	}
}
