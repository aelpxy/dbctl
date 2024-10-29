package cmd

import (
	"fmt"
	"time"

	"github.com/aelpxy/dbctl/docker"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:     "shell <container-id>",
	Short:   "Connect to a running database",
	Example: "dbctl shell container-id",
	Aliases: []string{"enter"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		connectToContainer(args[0])
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}

func connectToContainer(containerId string) {
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = fmt.Sprintf("Connecting to %s... \n", containerId)
	spinner.Color("green")
	spinner.Start()

	docker.ShellConnect(containerId)
}
