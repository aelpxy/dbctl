package cmd

import (
	"fmt"
	"log"

	"github.com/aelpxy/dbctl/docker"
	"github.com/spf13/cobra"
)

func init() {
	deleteCmd.Flags().BoolP("force", "f", true, "Delete the associated volume")

	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:     "delete <container-id>",
	Short:   "Stop and delete a database",
	Long:    "This command stops and removes a database container",
	Aliases: []string{"rm"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		containerId := args[0]

		deleteVolume, err := cmd.Flags().GetBool("force")

		if err != nil {
			log.Fatalf("Error getting volume flag: %v", err)
		}

		err = docker.DeleteContainer(containerId, deleteVolume)

		if err != nil {
			log.Fatalf("error deleting container: %v", err)
		}

		fmt.Printf("Container %s has been deleted.\n", containerId)

	},
}
