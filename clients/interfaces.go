package clients

import (
	"context"
	"io"

	"github.com/forta-protocol/forta-node/domain"

	"github.com/docker/docker/api/types"
	"github.com/golang/protobuf/proto"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/protocol"
)

// DockerClient is a client interface for interacting with docker
type DockerClient interface {
	PullImage(ctx context.Context, refStr string) error
	CreatePublicNetwork(ctx context.Context, name string) (string, error)
	CreateInternalNetwork(ctx context.Context, name string) (string, error)
	AttachNetwork(ctx context.Context, containerID string, networkID string) error
	GetContainers(ctx context.Context) (DockerContainerList, error)
	GetFortaServiceContainers(ctx context.Context) (fortaContainers DockerContainerList, err error)
	GetContainerByName(ctx context.Context, name string) (*types.Container, error)
	GetContainerByID(ctx context.Context, id string) (*types.Container, error)
	StartContainer(ctx context.Context, config DockerContainerConfig) (*DockerContainer, error)
	StopContainer(ctx context.Context, ID string) error
	InterruptContainer(ctx context.Context, ID string) error
	WaitContainerExit(ctx context.Context, id string) error
	WaitContainerStart(ctx context.Context, id string) error
	Prune(ctx context.Context) error
	WaitContainerPrune(ctx context.Context, id string) error
	Nuke(ctx context.Context) error
	HasLocalImage(ctx context.Context, ref string) bool
	EnsureLocalImage(ctx context.Context, name, ref string) error
}

// MessageClient receives and publishes messages.
type MessageClient interface {
	Subscribe(subject string, handler interface{})
	Publish(subject string, payload interface{})
	PublishProto(subject string, payload proto.Message)
}

// AgentClient makes the gRPC requests to evaluate block and txs and receive results.
type AgentClient interface {
	Dial(config.AgentConfig) error
	protocol.AgentClient
	io.Closer
}

// AlertAPIClient calls an http api on the analyzer to store alerts
type AlertAPIClient interface {
	PostBatch(batch *domain.AlertBatch, token string) error
}
