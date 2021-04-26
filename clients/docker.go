package clients

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

// DockerContainer is a resulting container reference, including the ID and configuration
type DockerContainer struct {
	ID     string
	Config DockerContainerConfig
}

// DockerContainerConfig is configuration for a particular container
type DockerContainerConfig struct {
	Name           string
	Image          string
	Env            map[string]string
	LinkNetworkIDs []string
	NetworkID      string
	Ports          map[string]string
}

// DockerClient is a client interface for interacting with docker
type DockerClient interface {
	CreatePublicNetwork(ctx context.Context, name string) (string, error)
	CreateInternalNetwork(ctx context.Context, name string) (string, error)
	AttachNetwork(ctx context.Context, containerID string, networkID string) error
	StartContainer(ctx context.Context, config DockerContainerConfig) (*DockerContainer, error)
	StopContainer(ctx context.Context, ID string) error
	Prune(ctx context.Context) error
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

func (d *dockerClient) Prune(ctx context.Context) error {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}
	zephyrFilter := filters.NewArgs(filters.Arg("label", "zephyr"))
	res, err := cli.NetworksPrune(ctx, zephyrFilter)
	if err != nil {
		return err
	}
	for _, nw := range res.NetworksDeleted {
		log.Infof("pruned network %s", nw)
	}

	cpRes, err := cli.ContainersPrune(ctx, zephyrFilter)
	if err != nil {
		return err
	}
	for _, cp := range cpRes.ContainersDeleted {
		log.Infof("pruned container %s", cp)
	}

	return nil
}

func (d *dockerClient) CreatePublicNetwork(ctx context.Context, name string) (string, error) {
	return d.createNetwork(ctx, name, false)
}

func (d *dockerClient) CreateInternalNetwork(ctx context.Context, name string) (string, error) {
	return d.createNetwork(ctx, name, true)
}

func (d *dockerClient) createNetwork(ctx context.Context, name string, internal bool) (string, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return "", err
	}
	resp, err := cli.NetworkCreate(ctx, name, types.NetworkCreate{
		Labels:   map[string]string{"zephyr": "true"},
		Internal: internal,
	})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (d *dockerClient) AttachNetwork(ctx context.Context, containerID string, networkID string) error {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}
	return cli.NetworkConnect(ctx, networkID, containerID, nil)
}

// StartContainer kicks off a container as a daemon and returns a summary of the container
func (d *dockerClient) StartContainer(ctx context.Context, config DockerContainerConfig) (*DockerContainer, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}

	cont, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:  config.Image,
			Env:    config.envVars(),
			Labels: map[string]string{"zephyr": "true"},
		},
		&container.HostConfig{
			NetworkMode: container.NetworkMode(config.NetworkID),
		}, nil, config.Name)

	if err != nil {
		return nil, err
	}

	if err := cli.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	for _, nwID := range config.LinkNetworkIDs {
		if err := d.AttachNetwork(ctx, cont.ID, nwID); err != nil {
			log.Error("error attaching network", err)
			return nil, err
		}
	}

	log.Infof("Container %s is started", cont.ID)
	return &DockerContainer{ID: cont.ID, Config: config}, nil
}

// StopContainer kills a container by ID
func (d *dockerClient) StopContainer(ctx context.Context, ID string) error {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}
	return cli.ContainerKill(ctx, ID, "SIGKILL")
}

// NewDockerClient creates a new docker client
func NewDockerClient() *dockerClient {
	return &dockerClient{}
}
