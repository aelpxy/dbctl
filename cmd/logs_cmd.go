package cmd

import (
	"github.com/aelpxy/dbctl/docker"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:     "logs <container-id>",
	Short:   "Stream live logs of a database",
	Example: "dbctl logs container-id",
	Aliases: []string{"tail"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		streamLogs(args[0])
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}

func streamLogs(containerId string) {
	docker.StreamLogs(containerId)
}
