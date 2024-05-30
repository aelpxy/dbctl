package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/aelpxy/dbctl/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var Ctx context.Context

func DockerClient() (*client.Client, error) {
	Ctx = context.Background()

	apiClient, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		return nil, err
	}

	_, err = apiClient.NetworkInspect(Ctx, config.DockerNetworkName, types.NetworkInspectOptions{})

	// basically checks if `config.DockerNetworkName` exists otherwise creates
	if err != nil {
		_, err = apiClient.NetworkCreate(Ctx, config.DockerNetworkName, types.NetworkCreate{})

		if err != nil {
			log.Fatalln(fmt.Errorf("error creating docker network: %w", err))
		}
	}

	return apiClient, nil
}
