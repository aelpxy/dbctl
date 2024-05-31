package cmd

import (
	"github.com/aelpxy/dbctl/docker"
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:     "shell <container-id>",
	Short:   "Connect to a running database",
	Long:    `Connect to a running database and open an interactive shell session inside it`,
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
	docker.ShellConnect(containerID)
}
