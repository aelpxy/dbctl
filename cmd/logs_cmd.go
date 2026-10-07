package cmd

import (
	"log"

	"github.com/aelpxy/dbctl/docker"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:     "logs <container-id>",
	Short:   "Stream live logs of a database",
	Example: "  dbctl logs container-id\n  dbctl logs container-id --tail 100 --follow=false",
	Aliases: []string{"tail"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		follow, _ := cmd.Flags().GetBool("follow")
		tail, _ := cmd.Flags().GetString("tail")

		if err := docker.StreamLogs(args[0], follow, tail); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	logsCmd.Flags().BoolP("follow", "f", true, "Follow log output.")
	logsCmd.Flags().StringP("tail", "t", "all", "Number of lines to show from the end of the logs.")

	rootCmd.AddCommand(logsCmd)
}
