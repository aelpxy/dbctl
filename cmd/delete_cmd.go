package cmd

import (
	"fmt"
	"log"

	"github.com/aelpxy/dbctl/docker"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <container-id>",
	Short:   "Stop and delete a database",
	Long:    "This command stops and removes a database container",
	Aliases: []string{"rm"},
	Example: `dbctl delete container-id
dbctl delete container-id --v true
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		containerId := args[0]

		deleteVolume, err := cmd.Flags().GetBool("f")

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

func init() {
	deleteCmd.Flags().BoolP("force", "f", true, "Delete the associated volume")

	rootCmd.AddCommand(deleteCmd)
}
