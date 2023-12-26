package utils

import (
	"github.com/docker/docker/client"
)

func InitializeDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}
