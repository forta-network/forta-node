package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	log "github.com/sirupsen/logrus"
)

// DockerContainer is a resulting container reference, including the ID and configuration
type DockerContainer struct {
	ID     string
	Config DockerContainerConfig
}

// DockerContainerConfig is configuration for a particular container
type DockerContainerConfig struct {
	Name  string
	Image string
	Env   map[string]string
}

// DockerClient is a client interface for interacting with docker
type DockerClient interface {
	StartContainer(ctx context.Context, config DockerContainerConfig) (DockerContainer, error)
	StopContainer(ID string) error
}

type dockerClient struct {
}

func (cfg DockerContainerConfig) envVars() []string {
	var results []string
	for k, v := range cfg.Env {
		results = append(results, fmt.Sprintf("%s=%s", k, v))
	}
	return results
}

// StartContainer kicks off a container as a daemon and returns a summary of the container
func (d *dockerClient) StartContainer(ctx context.Context, config DockerContainerConfig) (DockerContainer, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return DockerContainer{}, err
	}

	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "8000",
	}
	containerPort, err := nat.NewPort("tcp", "80")
	if err != nil {
		return DockerContainer{}, err
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	cont, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: config.Image,
			Env:   config.envVars(),
		},
		&container.HostConfig{
			NetworkMode:  "host", //TODO: update to inter-container awareness
			PortBindings: portBinding,
		}, nil, config.Name)

	if err != nil {
		return DockerContainer{}, err
	}

	if err := cli.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{}); err != nil {
		return DockerContainer{}, err
	}

	log.Infof("Container %s is started", cont.ID)
	return DockerContainer{ID: cont.ID, Config: config}, nil
}

// StopContainer kills a container by ID
func (d *dockerClient) StopContainer(ID string) error {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}
	timeout := time.Second * 30
	return cli.ContainerStop(context.Background(), ID, &timeout)
}

// NewDockerClient creates a new docker client
func NewDockerClient() *dockerClient {
	return &dockerClient{}
}
