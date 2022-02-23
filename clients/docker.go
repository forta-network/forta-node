package clients

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/utils/workers"
	log "github.com/sirupsen/logrus"
)

const dockerResourcesLabel = "network.forta"

var defaultLabels = map[string]string{dockerResourcesLabel: "true"}

// Client errors
var (
	ErrContainerNotFound = errors.New("container not found")
)

// DockerContainer is a resulting container reference, including the ID and configuration
type DockerContainer struct {
	Name      string
	ID        string
	ImageHash string
	Config    DockerContainerConfig
}

// DockerContainerConfig is configuration for a particular container
type DockerContainerConfig struct {
	Name            string
	Image           string
	Env             map[string]string
	LinkNetworkIDs  []string
	NetworkID       string
	Ports           map[string]string
	PublishAllPorts bool // auto-publishing ports EXPOSEd in Dockerfile
	Volumes         map[string]string
	Files           map[string][]byte
	MaxLogSize      string
	MaxLogFiles     int
	CPUQuota        int64
	Memory          int64
	Cmd             []string
	DialHost        bool
}

// DockerContainerList contains the full container data.
type DockerContainerList []types.Container

// FindByID finds the container by the ID.
func (dcl DockerContainerList) FindByID(id string) (*types.Container, bool) {
	for _, c := range dcl {
		if c.ID == id {
			return &c, true
		}
	}
	return nil, false
}

// FindByName finds the container by the name.
func (dcl DockerContainerList) FindByName(name string) (*types.Container, bool) {
	for _, c := range dcl {
		for _, n := range c.Names {
			if n == name || n == fmt.Sprintf("/%s", name) {
				return &c, true
			}
		}
	}
	return nil, false
}

type dockerClient struct {
	cli      *client.Client
	workers  *workers.Group
	username string
	password string
	labels   map[string]string
}

func (cfg DockerContainerConfig) envVars() []string {
	var results []string
	for k, v := range cfg.Env {
		results = append(results, fmt.Sprintf("%s=%s", k, v))
	}
	return results
}

func registryAuthValue(username, password string) string {
	if username == "" && password == "" {
		return ""
	}
	jsonBytes, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	return base64.StdEncoding.EncodeToString(jsonBytes)
}

// PullImage pulls an image using the given ref.
func (d *dockerClient) PullImage(ctx context.Context, refStr string) error {
	return d.workers.Execute(func() ([]interface{}, error) {
		return nil, d.pullImage(ctx, refStr)
	}).Error
}

func (d *dockerClient) pullImage(ctx context.Context, refStr string) error {
	r, err := d.cli.ImagePull(ctx, refStr, types.ImagePullOptions{
		RegistryAuth: registryAuthValue(d.username, d.password),
	})
	if err != nil {
		return err
	}
	defer r.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	respStr := strings.ToLower(string(b))
	if strings.Contains(respStr, "downloaded") || strings.Contains(respStr, "up to date") {
		return nil
	}
	return fmt.Errorf("unexpected image pull response: %s", string(b))
}

func (d *dockerClient) Prune(ctx context.Context) error {
	filter := d.labelFilter()
	res, err := d.cli.NetworksPrune(ctx, filter)
	if err != nil {
		return err
	}
	for _, nw := range res.NetworksDeleted {
		log.Infof("pruned network %s", nw)
	}

	cpRes, err := d.cli.ContainersPrune(ctx, filter)
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
	// Reuse if network exists.
	networks, err := d.cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return "", err
	}
	for _, network := range networks {
		if network.Name == name {
			return network.ID, nil
		}
	}

	resp, err := d.cli.NetworkCreate(ctx, name, types.NetworkCreate{
		Labels:   d.labels,
		Internal: internal,
	})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (d *dockerClient) AttachNetwork(ctx context.Context, containerID string, networkID string) error {
	err := d.cli.NetworkConnect(ctx, networkID, containerID, nil)
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "already exists") {
		return nil
	}
	return err
}

func withTcp(port string) string {
	return fmt.Sprintf("%s/tcp", port)
}

// copyFile copies content bytes into container at /filename
func copyFile(cli *client.Client, ctx context.Context, filename string, content []byte, containerId string) error {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	err := tw.WriteHeader(&tar.Header{
		Name: filename,
		Mode: 0400,
		Size: int64(len(content)),
	})
	if err != nil {
		return err
	}
	_, err = tw.Write(content)
	if err != nil {
		return err
	}
	err = tw.Close()
	if err != nil {
		return err
	}
	return cli.CopyToContainer(ctx, containerId, "/", &buf, types.CopyToContainerOptions{})
}

// GetContainers returns all of the containers.
func (d *dockerClient) GetContainers(ctx context.Context) (DockerContainerList, error) {
	return d.cli.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: d.labelFilter(),
	})
}

// GetFortaServiceContainers returns all of the non-agent forta containers.
func (d *dockerClient) GetFortaServiceContainers(ctx context.Context) (fortaContainers DockerContainerList, err error) {
	containers, err := d.cli.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: d.labelFilter(),
	})
	for _, container := range containers {
		if !strings.Contains(container.Names[0][1:], "forta-agent") {
			fortaContainers = append(fortaContainers, container)
		}
	}
	return
}

// GetContainerByName gets a container by using a name lookup over all containers.
func (d *dockerClient) GetContainerByName(ctx context.Context, name string) (*types.Container, error) {
	containers, err := d.GetContainers(ctx)
	if err != nil {
		return nil, err
	}
	for _, c := range containers {
		if c.Names[0][1:] == name {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("%w with name '%s'", ErrContainerNotFound, name)
}

// GetContainerByName gets a container by using an ID lookup over all containers.
func (d *dockerClient) GetContainerByID(ctx context.Context, id string) (*types.Container, error) {
	containers, err := d.GetContainers(ctx)
	if err != nil {
		return nil, err
	}
	for _, c := range containers {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("%w with id '%s'", ErrContainerNotFound, id)
}

// Nuke makes sure that all running Forta containers are stopped and pruned, quickly enough.
func (d *dockerClient) Nuke(ctx context.Context) error {
	var err error
	for i := 0; i < 4; i++ {
		err = d.nuke(ctx)
		if err == nil {
			return nil
		}
		log.WithError(err).Error("failed to nuke - retrying")
	}
	return fmt.Errorf("all nuke retries failed: %v", err)
}

func (d *dockerClient) nuke(ctx context.Context) error {
	// step 1: put the supervisor to the top of the list so it doesn't do funny restarts
	containers, err := d.GetContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get forta containers list: %v", err)
	}
	supervisorContainer, err := d.GetContainerByName(ctx, config.DockerSupervisorContainerName)
	if err == nil {
		containers = append([]types.Container{*supervisorContainer}, containers...)
	}
	if err != nil && !errors.Is(err, ErrContainerNotFound) {
		return fmt.Errorf("unexpected error while getting supervisor container: %v", err)
	}

	// step 1: stop all and wait until each exit
	for _, container := range containers {
		if err := d.StopContainer(ctx, container.ID); err != nil {
			return fmt.Errorf("failed to stop: %v", err)
		}
		if err := d.WaitContainerExit(ctx, container.ID); err != nil {
			return err
		}
	}

	// step 2: prune everything
	if err := d.Prune(ctx); err != nil {
		return fmt.Errorf("failed to prune: %v", err)
	}

	// step 3: ensure that the containers are really pruned
	for _, container := range containers {
		if err := d.WaitContainerPrune(ctx, container.ID); err != nil {
			return err
		}
	}

	return nil
}

// StartContainer kicks off a container as a daemon and returns a summary of the container
func (d *dockerClient) StartContainer(ctx context.Context, config DockerContainerConfig) (*DockerContainer, error) {
	log.WithFields(log.Fields{
		"image": config.Image,
		"name":  config.Name,
	}).Info("starting container")
	containers, err := d.GetContainers(ctx)
	if err != nil {
		return nil, err
	}
	// If we already have the container but it is not running, then just start it.
	var foundContainer *types.Container
	for _, c := range containers {
		if len(c.Names) == 0 {
			continue
		}
		foundName := c.Names[0][1:] // remove / in the beginning
		if foundName == config.Name {
			foundContainer = &c
			break
		}
	}
	if foundContainer != nil {
		if err := d.cli.ContainerStart(ctx, foundContainer.ID, types.ContainerStartOptions{}); err != nil {
			return nil, err
		}
		inspection, err := d.cli.ContainerInspect(ctx, foundContainer.ID)
		if err != nil {
			return nil, err
		}
		log.Infof("Container %s (%s) is started", foundContainer.ID, config.Name)
		return &DockerContainer{Name: config.Name, ID: foundContainer.ID, Config: config, ImageHash: inspection.Image}, nil
	}

	bindings := make(map[nat.Port][]nat.PortBinding)
	ps := make(nat.PortSet)
	for hp, cp := range config.Ports {
		hostIP := "0.0.0.0"
		parts := strings.Split(hp, ":")
		if len(parts) == 2 {
			hostIP = parts[0]
			hp = parts[1]
		}
		contPort := nat.Port(withTcp(cp))
		ps[contPort] = struct{}{}
		bindings[contPort] = []nat.PortBinding{{
			HostPort: hp,
			HostIP:   hostIP,
		}}
	}

	var volumes []string
	for hostVol, containerMnt := range config.Volumes {
		volumes = append(volumes, fmt.Sprintf("%s:%s", hostVol, containerMnt))
	}

	maxLogSize := config.MaxLogSize
	if maxLogSize == "" {
		maxLogSize = "10m"
	}

	maxLogFiles := config.MaxLogFiles
	if maxLogFiles == 0 {
		maxLogFiles = 10
	}

	cntCfg := &container.Config{
		Image:  config.Image,
		Env:    config.envVars(),
		Labels: d.labels,
	}

	if len(config.Cmd) > 0 {
		cntCfg.Cmd = config.Cmd
	}

	hostCfg := &container.HostConfig{
		NetworkMode:     container.NetworkMode(config.NetworkID),
		PortBindings:    bindings,
		PublishAllPorts: config.PublishAllPorts,
		Binds:           volumes,
		LogConfig: container.LogConfig{
			Config: map[string]string{
				"max-file": fmt.Sprintf("%d", maxLogFiles),
				"max-size": maxLogSize,
			},
			Type: "json-file",
		},
		Resources: container.Resources{
			CPUQuota: config.CPUQuota,
			Memory:   config.Memory,
		},
	}

	if config.DialHost {
		hostCfg.ExtraHosts = append(hostCfg.ExtraHosts, "host.docker.internal:host-gateway")
	}

	cont, err := d.cli.ContainerCreate(
		ctx,
		cntCfg,
		hostCfg, nil, config.Name)

	if err != nil {
		return nil, err
	}

	for fn, b := range config.Files {
		if err := copyFile(d.cli, ctx, fn, b, cont.ID); err != nil {
			return nil, err
		}
	}

	if err := d.cli.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	for _, nwID := range config.LinkNetworkIDs {
		if err := d.AttachNetwork(ctx, cont.ID, nwID); err != nil {
			log.Error("error attaching network", err)
			return nil, err
		}
	}

	inspection, err := d.cli.ContainerInspect(ctx, cont.ID)
	if err != nil {
		return nil, err
	}

	log.Infof("Container %s (%s) is started", cont.ID, config.Name)
	return &DockerContainer{Name: config.Name, ID: cont.ID, Config: config, ImageHash: inspection.Image}, nil
}

// StopContainer kills a container by ID
func (d *dockerClient) StopContainer(ctx context.Context, ID string) error {
	log.WithFields(log.Fields{
		"id": ID,
	}).Info("stop container (SIGKILL)")
	err := d.cli.ContainerKill(ctx, ID, "SIGKILL")
	if err == nil {
		return nil
	}
	if isNoSuchContainerErr(err) || isNotRunningErr(err) {
		return nil
	}
	return err
}

// InterruptContainer stops a container by sending an interrupt signal.
func (d *dockerClient) InterruptContainer(ctx context.Context, ID string) error {
	log.WithFields(log.Fields{
		"id": ID,
	}).Info("stop container (SIGINT)")
	err := d.cli.ContainerKill(ctx, ID, "SIGINT")
	if err == nil {
		return nil
	}
	if isNoSuchContainerErr(err) || isNotRunningErr(err) {
		return nil
	}
	return err
}

func isNoSuchContainerErr(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "no such container")
}

func isNotRunningErr(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "is not running")
}

// WaitContainerExit waits for container exit by checking periodically.
func (d *dockerClient) WaitContainerExit(ctx context.Context, id string) error {
	ticker := time.NewTicker(time.Second)
	logger := log.WithFields(log.Fields{
		"id": id,
	})

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	for range ticker.C {
		logger.Info("waiting for container exit")
		c, err := d.GetContainerByID(ctx, id)
		if err != nil && errors.Is(err, ErrContainerNotFound) {
			logger.Info("no need to wait for container exit - not found")
			return nil
		}
		if err != nil {
			logger.WithError(err).Error("failed while waiting for container exit")
			return err
		}
		if c.State == "exited" || c.State == "created" {
			return nil
		}
		logger.WithField("containerState", c.State).Info("still waiting for exit")
	}
	return nil
}

// WaitContainerStart waits for container start by checking periodically.
func (d *dockerClient) WaitContainerStart(ctx context.Context, id string) error {
	ticker := time.NewTicker(time.Second)
	start := time.Now()
	logger := log.WithFields(log.Fields{
		"id": id,
	})

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	for t := range ticker.C {
		logger.Info("waiting for container start")
		c, err := d.GetContainerByID(ctx, id)
		if err == nil && c != nil && c.State == "running" {
			return nil
		}
		if err != nil {
			return err
		}
		// if the conditions are not met within 30 seconds, it's a failure
		if t.After(start.Add(time.Second * 30)) {
			return errors.New("container did not start")
		}
	}
	return nil
}

// WaitContainerPrune waits for container prune by checking periodically.
func (d *dockerClient) WaitContainerPrune(ctx context.Context, id string) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	logger := log.WithFields(log.Fields{
		"id": id,
	})

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	for {
		logger.Infof("waiting for container prune")
		c, err := d.GetContainerByID(ctx, id)
		if err != nil && errors.Is(err, ErrContainerNotFound) {
			return nil
		}
		if err != nil {
			logger.WithError(err).Error("error while waiting for prune")
			return err
		}
		logger.WithField("containerState", c.State).Info("container state while waiting for prune")
		if !(c.State == "exited" || c.State == "dead") {
			err = fmt.Errorf("cannot prune container with status '%s' - container needs to stop first", c.State)
			logger.WithError(err).Error("error while waiting for prune")
			return err
		}
		<-ticker.C
	}
}

// HasLocalImage checks if we have an image locally.
func (d *dockerClient) HasLocalImage(ctx context.Context, ref string) bool {
	_, _, err := d.cli.ImageInspectWithRaw(ctx, ref)
	return err == nil
}

// EnsureLocalImage ensures that we have the image locally.
func (d *dockerClient) EnsureLocalImage(ctx context.Context, name, ref string) error {
	log.WithFields(log.Fields{
		"image": ref,
		"name":  name,
	}).Info("ensuring local image")
	if d.HasLocalImage(ctx, ref) {
		log.Infof("found local image for '%s': %s", name, ref)
		return nil
	}
	err := d.PullImage(ctx, ref)
	if err != nil {
		return fmt.Errorf("failed to pull image (%s, %s): %v", name, ref, err)
	}
	log.Infof("pulled image for '%s': %s", name, ref)
	return nil
}

func (d *dockerClient) labelFilter() filters.Args {
	filter := filters.NewArgs()
	for k, v := range d.labels {
		filter.Add("label", fmt.Sprintf("%s=%s", k, v))
	}
	return filter
}

func initLabels(name string) map[string]string {
	result := make(map[string]string)
	for k, v := range defaultLabels {
		result[k] = v
	}
	if len(name) == 0 {
		return result
	}
	result["network.forta.supervisor"] = name
	return result
}

// NewDockerClient creates a new docker client
func NewDockerClient(name string) (*dockerClient, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	return &dockerClient{
		cli:     cli,
		workers: workers.New(10),
		labels:  initLabels(name),
	}, nil
}

// NewAuthDockerClient creates a new docker client with credentials
func NewAuthDockerClient(name string, username, password string) (*dockerClient, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	return &dockerClient{
		cli:      cli,
		workers:  workers.New(10),
		username: username,
		password: password,
		labels:   initLabels(name),
	}, nil
}
