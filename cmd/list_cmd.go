package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/aelpxy/dbctl/config"
	"github.com/aelpxy/dbctl/docker"
	"github.com/aelpxy/dbctl/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all running databases",
	Long:  "List all running databases managed by dbctl",
	Run: func(cmd *cobra.Command, args []string) {
		listDatabases()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listDatabases() {
	dockerClient, err := docker.DockerClient()
	if err != nil {
		log.Fatalf("error creating Docker client: %v", err)
	}

	containers, err := dockerClient.ContainerList(docker.Ctx, container.ListOptions{All: true})

	if err != nil {
		log.Fatalf("error listing containers: %v", err)
	}

	var filteredContainers []types.Container

	for _, container := range containers {
		// this filters the containers and only finds ones starting with dbctl.*
		if strings.HasPrefix(strings.TrimPrefix(container.Names[0], "/"), config.DockerContainerPrefix) {
			filteredContainers = append(filteredContainers, container)
		}
	}

	if len(filteredContainers) == 0 {
		log.Println("No databases are currently running.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name", "Image", "Status"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, container := range filteredContainers {
		// note: i don't know why docker has a leading slash
		// this is pretty long i may take a look in the future
		table.Rich([]string{container.ID[:12], strings.TrimPrefix(container.Names[0], "/"), container.Image, container.State}, []tablewriter.Colors{{tablewriter.FgGreenColor, tablewriter.Bold}, {tablewriter.FgBlueColor, tablewriter.Bold}, {tablewriter.FgGreenColor}, utils.GetColorBasedOnStatus(container.State)})
	}

	table.Render()
}
