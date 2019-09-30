package docker

import (
	"context"
	"errors"
	"fmt"

	docker "github.com/docker/docker/client"
)

// NewClient returns a new Docker client
func NewClient() (*docker.Client, error) {
	dockerClient, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error creating docker client. %v", err))
	}

	return dockerClient, nil
}

// IsDockerAvailable performs an 'info' call to the docker server to see if it is available
// Returns true if it is, otherwise returns false and the error message
func IsDockerAvailable() (bool, error) {

	dockerClient, err := NewClient()
	if err != nil {
		return false, err
	}

	// Make a call to the server to see if its reachable
	ctx := context.Background()
	_, err = dockerClient.Info(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}
