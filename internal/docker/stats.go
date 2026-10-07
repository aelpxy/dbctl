package docker

import (
	"context"
	"encoding/json/v2"
	"fmt"

	"github.com/docker/docker/api/types/container"
)

const (
	bytesPerMB = 1 << 20
	percent    = 100
)

// Stats is a resource usage snapshot of a database container.
type Stats struct {
	CPUPercent  float64
	MemoryUsage uint64
	MemoryLimit uint64
}

// MemoryUsageMB returns the memory usage in megabytes.
func (s Stats) MemoryUsageMB() float64 {
	return float64(s.MemoryUsage) / bytesPerMB
}

// MemoryLimitMB returns the memory limit in megabytes.
func (s Stats) MemoryLimitMB() float64 {
	return float64(s.MemoryLimit) / bytesPerMB
}

// Stats returns a resource usage snapshot of the database container.
func (c *Client) Stats(ctx context.Context, id string) (Stats, error) {
	resp, err := c.api.ContainerStats(ctx, id, false)
	if err != nil {
		return Stats{}, fmt.Errorf("get container stats %s: %w", id, err)
	}

	defer func() { _ = resp.Body.Close() }()

	var s container.StatsResponse

	if err := json.UnmarshalRead(resp.Body, &s); err != nil {
		return Stats{}, fmt.Errorf("decode container stats: %w", err)
	}

	return Stats{
		CPUPercent:  cpuPercent(&s),
		MemoryUsage: s.MemoryStats.Usage,
		MemoryLimit: s.MemoryStats.Limit,
	}, nil
}

func cpuPercent(s *container.StatsResponse) float64 {
	cpuDelta := float64(s.CPUStats.CPUUsage.TotalUsage) - float64(s.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(s.CPUStats.SystemUsage) - float64(s.PreCPUStats.SystemUsage)

	if cpuDelta <= 0 || systemDelta <= 0 {
		return 0
	}

	// PercpuUsage is empty on cgroup v2 hosts
	cpus := float64(s.CPUStats.OnlineCPUs)
	if cpus == 0 {
		cpus = float64(len(s.CPUStats.CPUUsage.PercpuUsage))
	}

	return cpuDelta / systemDelta * cpus * percent
}
