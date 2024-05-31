package utils

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
	"github.com/olekukonko/tablewriter"
)

func GeneratePassword(length int) string {
	// not the most secure way but does the job
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func GetAvailablePort() int {
	// not the best way
	for i := 0; i < 10; i++ {
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

func IsDockerInstalled() bool {
	_, err := exec.Command("docker", "version").Output()
	return err == nil
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
	parts := strings.Split(part, ":")
	dbType := parts[0]
	imageVersion := ""

	if len(parts) > 1 {
		imageVersion = parts[1]
	}

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
