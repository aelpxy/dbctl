package utils

import (
	crand "crypto/rand"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"net"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
	"github.com/olekukonko/tablewriter"
)

func GeneratePassword(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	for i := range b {
		n, err := crand.Int(crand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			log.Fatalf("error generating password: %v", err)
		}

		b[i] = charset[n.Int64()]
	}

	return string(b)
}

func GetAvailablePort() int {
	// not the best way
	for range 10 {
		port := 30000 + rand.Intn(10000)
		if isPortAvailable(port) {
			return port
		}
	}

	log.Fatalf("unable to find an available port to use")

	return 0
}

// checks if the port is available
func isPortAvailable(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	conn, err := net.Listen("tcp", addr)

	if err != nil {
		return false
	}

	conn.Close()

	return true
}

func GetColorBasedOnStatus(status string) tablewriter.Colors {
	switch status {
	case "running":
		return tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.Bold}
	case "exited":
		return tablewriter.Colors{tablewriter.FgRedColor, tablewriter.Bold}
	case "paused":
		return tablewriter.Colors{tablewriter.FgYellowColor, tablewriter.Bold}
	default:
		return tablewriter.Colors{tablewriter.FgWhiteColor, tablewriter.BgBlackColor}
	}
}

func ParseDBTypeAndVersion(part string) (string, string) {
	dbType, imageVersion, _ := strings.Cut(part, ":")

	return dbType, imageVersion
}

func FormatPorts(ports nat.PortMap) string {
	var formattedPorts []string
	for port, bindings := range ports {
		for _, binding := range bindings {
			formattedPorts = append(formattedPorts, fmt.Sprintf("%s:%s->%s", binding.HostIP, binding.HostPort, port.Port()))
		}
	}
	return strings.Join(formattedPorts, ", ")
}

func FormatVolumes(mounts []types.MountPoint) string {
	var formattedVolumes []string
	for _, mount := range mounts {
		formattedVolumes = append(formattedVolumes, mount.Name)
	}
	return strings.Join(formattedVolumes, ", ")
}

func FormatStorage(usage, limit uint64) string {
	return fmt.Sprintf("%.2f GB / %.2f GB", float64(usage)/1024/1024/1024, float64(limit)/1024/1024/1024)
}
