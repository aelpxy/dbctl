package docker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
)

type ContainerStats struct {
	CPUPercentage float64
	MemoryUsage   uint64
	MemoryLimit   uint64
	StorageUsage  uint64
	StorageLimit  uint64
}

func GetContainerStats(containerId string) (*ContainerStats, error) {
	dockerClient, err := DockerClient()

	if err != nil {
		return nil, fmt.Errorf("error creating docker client: %v", err)
	}

	stats, err := dockerClient.ContainerStats(context.Background(), containerId, false)

	if err != nil {
		return nil, fmt.Errorf("error getting container stats: %v", err)
	}
	defer stats.Body.Close()

	var containerStats types.StatsJSON

	err = json.NewDecoder(stats.Body).Decode(&containerStats)

	if err != nil {
		return nil, fmt.Errorf("error decoding container stats: %v", err)
	}

	cpuDelta := float64(containerStats.CPUStats.CPUUsage.TotalUsage - containerStats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(containerStats.CPUStats.SystemUsage - containerStats.PreCPUStats.SystemUsage)
	cpuPercentage := (cpuDelta / systemDelta) * float64(len(containerStats.CPUStats.CPUUsage.PercpuUsage)) * 100

	return &ContainerStats{
		CPUPercentage: cpuPercentage,
		MemoryUsage:   containerStats.MemoryStats.Usage,
		MemoryLimit:   containerStats.MemoryStats.Limit,
		StorageUsage:  containerStats.StorageStats.ReadSizeBytes + containerStats.StorageStats.WriteSizeBytes,
		StorageLimit:  containerStats.StorageStats.ReadSizeBytes + containerStats.StorageStats.WriteSizeBytes,
	}, nil
}
