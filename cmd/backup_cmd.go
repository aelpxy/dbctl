package cmd

import (
	"github.com/aelpxy/dbctl/docker"
	"github.com/spf13/cobra"
)

// TODO: add support for multiple databases
var backupCmd = &cobra.Command{
	Use:     "backup <container-id>",
	Short:   "Backup a database",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"cp"},
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")

		docker.BackupContainer(args[0], output)
	},
}

func init() {
	backupCmd.Flags().StringP("output", "o", "", "Specify a output path for the database backup.")

	rootCmd.AddCommand(backupCmd)
}
