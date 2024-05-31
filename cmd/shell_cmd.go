package cmd

import (
	"log"

	"github.com/aelpxy/dbctl/docker"
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:     "shell <container-id>",
	Short:   "Connect to a running container",
	Long:    `Connect to a running container and open an interactive shell session inside it`,
	Example: "dbctl shell container-id",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		containerID := args[0]
		connectToContainer(containerID)
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}

func connectToContainer(containerID string) {
	err := docker.ShellConnect(containerID)

	if err != nil {
		log.Fatalf("error connecting to container: %v", err)
	}
}
