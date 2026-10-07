// Package hostnet looks up host network details for publishing database ports.
package hostnet

import (
	"context"
	"errors"
	"fmt"
	"net"
)

const (
	resolverAddress = "9.9.9.9:80"
	loopback        = "127.0.0.1"
)

// OutboundIP returns the address of the interface used for outbound traffic, or loopback when offline.
func OutboundIP(ctx context.Context) net.IP {
	d := net.Dialer{}

	// dialing udp sends no packets, it only resolves the outbound interface address
	conn, err := d.DialContext(ctx, "udp", resolverAddress)
	if err != nil {
		return net.ParseIP(loopback)
	}

	defer func() { _ = conn.Close() }()

	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return net.ParseIP(loopback)
	}

	return addr.IP
}

// FreePort returns a tcp port that is currently free on the host.
func FreePort(ctx context.Context) (int, error) {
	lc := net.ListenConfig{}

	l, err := lc.Listen(ctx, "tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("listen on a free port: %w", err)
	}

	defer func() { _ = l.Close() }()

	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		return 0, errors.New("listener has no tcp address")
	}

	return addr.Port, nil
}
