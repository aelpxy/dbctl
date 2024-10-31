package cmd

import (
	"fmt"
	"os"

	"github.com/aelpxy/dbctl/docker"
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:     "shell <container-id>",
	Short:   "Connect to a running database container",
	Example: "dbctl shell container-id",
	Aliases: []string{"enter", "sh", "connect"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return connectToContainer(args[0])
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}

func connectToContainer(containerID string) error {
	if err := docker.ShellConnect(containerID); err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to container: %v\n", err)
		return err
	}

	return nil
}
