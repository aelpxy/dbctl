package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aelpxy/dbctl/config"
	"github.com/aelpxy/dbctl/docker"
	"github.com/aelpxy/dbctl/structs"
)

type RetrieveDatabaseInfo struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Image        string            `json:"image"`
	Status       string            `json:"status"`
	CPUUsage     float64           `json:"cpu_usage"`
	MemoryUsage  MemoryStats       `json:"memory_usage"`
	StorageUsage StorageStats      `json:"storage_usage"`
	Ports        map[string]string `json:"ports,omitempty"`
	Volumes      []VolumeInfo      `json:"volumes,omitempty"`
}

type MemoryStats struct {
	Used  float64 `json:"used_mb"`
	Total float64 `json:"total_mb"`
}

type StorageStats struct {
	Used  uint64 `json:"used_bytes"`
	Total uint64 `json:"total_bytes"`
}

type VolumeInfo struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func RetrieveDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	containerId := r.PathValue("id")

	if containerId == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(structs.Response{
			Success: false,
			Message: "You need to provide the id to retrieve information about a specific database.",
		})

		return
	}

	containerInfo, err := docker.InspectContainer(containerId)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(structs.Response{
			Success: false,
			Message: "Error inspecting container.",
		})

		return
	}

	if !strings.HasPrefix(strings.TrimPrefix(containerInfo.Name, "/"), config.DockerContainerPrefix) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(structs.Response{
			Success: false,
			Message: "Your database is not managed by dbctl.",
		})

		return
	}

	stats, err := docker.GetContainerStats(containerId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(structs.Response{
			Success: false,
			Message: "Error getting container stats.",
		})

		return
	}

	ports := make(map[string]string)
	if containerInfo.NetworkSettings != nil {
		for port, bindings := range containerInfo.NetworkSettings.Ports {
			if len(bindings) > 0 {
				ports[string(port)] = bindings[0].HostPort
			}
		}
	}

	volumes := make([]VolumeInfo, 0)
	if containerInfo.Mounts != nil {
		for _, mount := range containerInfo.Mounts {
			volumes = append(volumes, VolumeInfo{
				Source: mount.Source,
				Target: mount.Destination,
			})
		}
	}

	dbInfo := RetrieveDatabaseInfo{
		ID:     containerInfo.ID[:12],
		Name:   strings.TrimPrefix(containerInfo.Name, "/"),
		Image:  containerInfo.Config.Image,
		Status: containerInfo.State.Status,
		MemoryUsage: MemoryStats{
			Used:  float64(stats.MemoryUsage) / 1024 / 1024,
			Total: float64(stats.MemoryLimit) / 1024 / 1024,
		},
		CPUUsage: stats.CPUPercentage,
		StorageUsage: StorageStats{
			Used:  stats.StorageUsage,
			Total: stats.StorageLimit,
		},
		Ports:   ports,
		Volumes: volumes,
	}

	response := structs.Response{
		Success: true,
		Data:    []any{dbInfo},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
