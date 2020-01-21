package docker

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

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

// IsInDocker return true if currently running in a container, false otherwise
// No clear consensus on how to do this so it may need to change in the future
func IsInDocker() bool {

	// Check if a .dockerenv exists
	// This should only be in docker containers but there's always a chance...
	_, err := os.Stat("/.dockerenv")
	if err == nil {
		return true
	}

	// Check if our current process is in a docker cgroup
	_, err = os.Stat("/proc/self/cgroup")
	if err == nil {
		if f, err := os.Open("/proc/self/cgroup"); err != nil {
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if strings.Contains(scanner.Text(), "docker") {
					return true
				}
			}
		}
	}

	return false
}
