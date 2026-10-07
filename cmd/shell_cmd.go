package cmd

import (
	"github.com/aelpxy/dbctl/docker"
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:     "shell <container-id>",
	Short:   "Connect to a running database container",
	Example: "  dbctl shell container-id",
	Aliases: []string{"enter", "sh", "connect"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return docker.ShellConnect(args[0])
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
