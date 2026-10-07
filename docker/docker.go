package docker

import (
	"context"
	"fmt"
	"sync"

	"github.com/aelpxy/dbctl/config"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

var Ctx = context.Background()

var (
	cachedClient *client.Client
	clientMu     sync.Mutex
)

func DockerClient() (*client.Client, error) {
	clientMu.Lock()
	defer clientMu.Unlock()

	if cachedClient != nil {
		return cachedClient, nil
	}

	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	if _, err := apiClient.Ping(Ctx); err != nil {
		apiClient.Close()
		return nil, fmt.Errorf("unable to reach the docker daemon, make sure docker is installed and running: %w", err)
	}

	// creates `config.DockerNetworkName` if it does not exist yet
	if _, err := apiClient.NetworkInspect(Ctx, config.DockerNetworkName, network.InspectOptions{}); err != nil {
		if _, err := apiClient.NetworkCreate(Ctx, config.DockerNetworkName, network.CreateOptions{}); err != nil {
			apiClient.Close()
			return nil, fmt.Errorf("error creating docker network: %w", err)
		}
	}

	cachedClient = apiClient

	return cachedClient, nil
}
