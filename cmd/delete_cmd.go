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
	Use:     "delete <container-id>...",
	Short:   "Stop and delete one or more databases",
	Long:    "This command stops and removes database containers along with their volumes (use --force=false to keep the volumes)",
	Aliases: []string{"rm"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteVolume, err := cmd.Flags().GetBool("force")
		if err != nil {
			log.Fatalf("error getting force flag: %v", err)
		}

		for _, containerId := range args {
			if err := docker.DeleteContainer(containerId, deleteVolume); err != nil {
				log.Fatalf("error deleting container %s: %v", containerId, err)
			}

			fmt.Printf("Container %s has been deleted.\n", containerId)
		}
	},
}
