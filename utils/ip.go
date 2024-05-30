package utils

import (
	"log"
	"net"

	"github.com/aelpxy/dbctl/config"
)

func GetIP() net.IP {
	conn, err := net.Dial("udp", config.DNSResolverAddress)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP
}
