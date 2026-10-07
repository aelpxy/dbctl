package cmd

import (
	"fmt"
	"os"

	"github.com/aelpxy/dbctl/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          config.CmdName,
	Short:        config.CmdShortDescription,
	Long:         config.CmdLongDescription,
	SuggestFor:   []string{"db"},
	SilenceUsage: true,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the current dbctl version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Version)
	},
}

func Execute() {
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
