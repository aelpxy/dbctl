package utils

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os/exec"

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