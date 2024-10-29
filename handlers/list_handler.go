package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aelpxy/dbctl/config"
	"github.com/aelpxy/dbctl/docker"
	"github.com/aelpxy/dbctl/structs"
	"github.com/docker/docker/api/types/container"
)

type ListDatabaseInfo struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Status  string            `json:"status"`
	State   string            `json:"state"`
	Image   string            `json:"image"`
	Created int64             `json:"created"`
	Labels  map[string]string `json:"labels"`
}

func ListDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	dockerClient, err := docker.DockerClient()

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(structs.Response{
			Success: false,
			Message: "Failed to connect to Docker.",
		})

		return
	}

	containers, err := dockerClient.ContainerList(docker.Ctx, container.ListOptions{All: true})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(structs.Response{
			Success: false,
			Message: "Failed to list containers.",
		})
		return
	}

	var filteredDatabases []ListDatabaseInfo

	for _, container := range containers {
		name := strings.TrimPrefix(container.Names[0], "/")

		if strings.HasPrefix(name, config.DockerContainerPrefix) {
			dbInfo := ListDatabaseInfo{
				ID:      container.ID,
				Name:    name,
				Status:  container.Status,
				State:   container.State,
				Image:   container.Image,
				Created: container.Created,
				Labels:  container.Labels,
			}
			filteredDatabases = append(filteredDatabases, dbInfo)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	data := make([]any, len(filteredDatabases))

	for i, db := range filteredDatabases {
		data[i] = db
	}

	response := structs.Response{
		Success: true,
		Data:    data,
	}

	if len(filteredDatabases) == 0 {
		response.Message = "No databases are currently running."
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(structs.Response{
			Success: false,
			Message: "Failed to encode response.",
		})

		return
	}
}
