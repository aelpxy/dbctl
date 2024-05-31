package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aelpxy/dbctl/config"
	"github.com/aelpxy/dbctl/docker"
	"github.com/aelpxy/dbctl/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <container-id>",
	Short: "Inspect a running database",
	Long:  "Inspect a running database managed by dbctl",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inspectDatabase(args[0])
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}

func inspectDatabase(containerId string) {
	containerInfo, err := docker.InspectContainer(containerId)

	if err != nil {
		log.Fatalf("error inspecting container: %v", err)
	}

	if !strings.HasPrefix(strings.TrimPrefix(containerInfo.Name, "/"), config.DockerContainerPrefix) {
		log.Fatalf("this container %s is not managed by dbctl", containerInfo.Name)
	}

	stats, err := docker.GetContainerStats(containerId)
	if err != nil {
		log.Fatalf("error getting container stats: %v", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.FgHiGreenColor, tablewriter.Bold},
	)

	table.Append([]string{"ID", containerInfo.ID[:12]})
	table.Append([]string{"Name", strings.TrimPrefix(containerInfo.Name, "/")})
	table.Append([]string{"Image", containerInfo.Config.Image})
	table.Append([]string{"Status", containerInfo.State.Status})
	table.Append([]string{"CPU Usage", fmt.Sprintf("%.2f%%", stats.CPUPercentage)})
	table.Append([]string{"Memory Usage", fmt.Sprintf("%.2f MB / %.2f MB", float64(stats.MemoryUsage)/1024/1024, float64(stats.MemoryLimit)/1024/1024)})
	// not sure if this actually works
	table.Append([]string{"Storage Usage", utils.FormatStorage(stats.StorageUsage, stats.StorageLimit)})

	if containerInfo.NetworkSettings != nil {
		table.Append([]string{"Ports", utils.FormatPorts(containerInfo.NetworkSettings.Ports)})
	} else {
		table.Append([]string{"Ports", "none"})
	}

	if containerInfo.Mounts != nil {
		table.Append([]string{"Volumes", utils.FormatVolumes(containerInfo.Mounts)})
	} else {
		table.Append([]string{"Volumes", "none"})
	}

	table.Render()
}
