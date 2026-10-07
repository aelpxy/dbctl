package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/aelpxy/dbctl/internal/docker"
	"github.com/aelpxy/dbctl/internal/proxy"
	"github.com/spf13/cobra"
)

const loopbackHost = "127.0.0.1"

var errNoLocalAddr = errors.New("listener has no tcp address")

func (a *app) newConnectCmd() *cobra.Command {
	port := 0

	cmd := &cobra.Command{
		Use:   "connect [name-or-id]",
		Short: "Open a local port that forwards to a database",
		Long: "Listen on 127.0.0.1 and forward every connection to the database, printing a ready-to-use " +
			"connection string. Runs until Ctrl-C.",
		Example: "  dbctl connect misty-river-bold-pine\n  dbctl connect misty-river-bold-pine --port 15432",
		Args:    cobra.MaximumNArgs(1),
		RunE: a.withDocker(func(ctx context.Context, c *docker.Client, args []string) error {
			id, err := databaseArg(ctx, c, args, "Connect to which database?")
			if err != nil {
				return err
			}

			return a.runConnect(ctx, c, id, port)
		}),
	}

	cmd.Flags().IntVarP(&port, "port", "p", 0, "Local port to listen on (default: the database's usual port, or a free one).")

	return cmd
}

func (a *app) runConnect(ctx context.Context, c *docker.Client, id string, port int) error {
	db, err := c.Inspect(ctx, id)
	if err != nil {
		return fmt.Errorf("inspect database: %w", err)
	}

	if !db.Running {
		return fmt.Errorf("connect to %s: %w", db.Name, docker.ErrNotRunning)
	}

	def, err := a.databases.Get(db.Type)
	if err != nil {
		return fmt.Errorf("resolve database type: %w", err)
	}

	published, err := connectionFor(def, db)
	if err != nil {
		return fmt.Errorf("find database address: %w", err)
	}

	l, err := listenLocal(ctx, port, def.Port())
	if err != nil {
		return err
	}

	local, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		return errNoLocalAddr
	}

	conn, err := newConnection(def, db.Name, loopbackHost, local.Port, published.Password)
	if err != nil {
		return err
	}

	target := net.JoinHostPort(cmp.Or(published.Host, loopbackHost), strconv.Itoa(published.Port))
	printForwarding(conn, target)

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))

	if err := proxy.Serve(ctx, l, target, logger); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("forward connections: %w", err)
	}

	printErr(successLine("Stopped forwarding to " + bold(db.Name)))

	return nil
}

// listenLocal prefers the requested port, then the database's usual port, then any free port.
func listenLocal(ctx context.Context, port, usual int) (net.Listener, error) {
	lc := net.ListenConfig{}

	if port != 0 {
		l, err := lc.Listen(ctx, "tcp", net.JoinHostPort(loopbackHost, strconv.Itoa(port)))
		if err != nil {
			return nil, fmt.Errorf("listen on port %d: %w", port, err)
		}

		return l, nil
	}

	if l, err := lc.Listen(ctx, "tcp", net.JoinHostPort(loopbackHost, strconv.Itoa(usual))); err == nil {
		return l, nil
	}

	l, err := lc.Listen(ctx, "tcp", net.JoinHostPort(loopbackHost, "0"))
	if err != nil {
		return nil, fmt.Errorf("listen on a free port: %w", err)
	}

	return l, nil
}

func printForwarding(conn *connectionJSON, target string) {
	body := strings.Join([]string{
		accent(conn.Name) + " is available locally",
		"",
		keyValues([][2]string{
			{"Local", conn.Host + ":" + strconv.Itoa(conn.Port)},
			{"Forwarding to", target},
		}),
		"",
		muted("Connection string"),
		fg(colorAccent).Render(conn.URL),
	}, "\n")

	printOut(panel(body))
	printOut(hintLine("Press Ctrl-C to stop."))
}
