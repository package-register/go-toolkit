package docker

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// DockerClient represents a Docker client with synchronization support.
type DockerClient struct {
	client *client.Client
	ctx    context.Context
	mu     sync.Mutex // Ensures thread safety for client operations
}

// NewDockerClient creates and returns a new Docker client.
func NewDockerClient() (*DockerClient, error) {
	dockerURL := os.Getenv("DOCKER_HOST") // From environment variable
	if dockerURL == "" {
		dockerURL = "unix:///var/run/docker.sock" // Default Unix Socket
	}

	cli, err := client.NewClientWithOpts(client.WithHost(dockerURL), client.WithVersion("1.41"))
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %v", err)
	}

	return &DockerClient{client: cli, ctx: context.Background()}, nil
}

// CreateContainer creates and starts a new Docker container.
func (dc *DockerClient) CreateContainer(image string, portMappings map[string]string, envVars map[string]string, volumes map[string]string) (string, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	// Timeout for Docker container creation and start
	timeout := 5 * time.Minute
	ctx, cancel := context.WithTimeout(dc.ctx, timeout)
	defer cancel()

	containerConfig := &container.Config{
		Image: image,
		Env:   formatEnvVars(envVars),
	}

	hostConfig := &container.HostConfig{
		PortBindings: formatPortBindings(portMappings),
		Binds:        formatVolumes(volumes),
	}

	containerResp, err := dc.client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to create container: %v", err)
	}

	err = dc.client.ContainerStart(ctx, containerResp.ID, container.StartOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to start container: %v", err)
	}

	return containerResp.ID, nil
}

// StopContainer stops a running Docker container by its ID.
func (dc *DockerClient) StopContainer(containerID string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	ctx, cancel := context.WithTimeout(dc.ctx, 30*time.Second)
	defer cancel()

	err := dc.client.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("failed to stop container %s: %v", containerID, err)
	}

	return nil
}

// RemoveContainer removes a stopped Docker container by its ID.
func (dc *DockerClient) RemoveContainer(containerID string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	ctx, cancel := context.WithTimeout(dc.ctx, 30*time.Second)
	defer cancel()

	err := dc.client.ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force:         true,
		RemoveLinks:   true,
		RemoveVolumes: true,
	})
	if err != nil {
		return fmt.Errorf("failed to remove container %s: %v", containerID, err)
	}

	return nil
}

// formatEnvVars formats environment variables as Docker requires.
func formatEnvVars(envVars map[string]string) []string {
	var result []string
	for key, value := range envVars {
		result = append(result, key+"="+value)
	}
	return result
}

// formatPortBindings formats port mappings as Docker requires.
func formatPortBindings(portMappings map[string]string) nat.PortMap {
	portMap := make(nat.PortMap)
	for port, mapping := range portMappings {
		portMap[nat.Port(port)] = []nat.PortBinding{
			{HostPort: mapping},
		}
	}
	return portMap
}

// formatVolumes formats volume mappings as Docker requires.
func formatVolumes(volumes map[string]string) []string {
	var result []string
	for src, dst := range volumes {
		result = append(result, src+":"+dst)
	}
	return result
}
