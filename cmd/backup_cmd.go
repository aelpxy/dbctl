package cmd

import (
	"fmt"
	"log"

	"github.com/aelpxy/dbctl/docker"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:     "backup <container-id>",
	Short:   "Backup a database (postgres, mysql, mariadb, mongo)",
	Example: "  dbctl backup container-id\n  dbctl backup container-id -o backup.sql",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"cp"},
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")

		path, err := docker.BackupContainer(args[0], output)
		if err != nil {
			log.Fatalf("error backing up database: %v", err)
		}

		fmt.Printf("Backup saved to %s\n", path)
	},
}

func init() {
	backupCmd.Flags().StringP("output", "o", "", "Specify an output path for the database backup (default: <name>-<timestamp>.<ext>).")

	rootCmd.AddCommand(backupCmd)
}
