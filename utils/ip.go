package utils

import (
	"net"

	"github.com/aelpxy/dbctl/config"
)

// no packets are sent, dialing udp only resolves the outbound interface address
func GetIP() net.IP {
	conn, err := net.Dial("udp", config.DNSResolverAddress)
	if err != nil {
		return net.IPv4(127, 0, 0, 1)
	}

	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP
}
