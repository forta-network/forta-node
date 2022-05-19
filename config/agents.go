package config

import (
	"strings"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
)

const (
	AgentGrpcPort = "50051"
)

type AgentConfig struct {
	ID         string  `yaml:"id" json:"id"`
	Image      string  `yaml:"image" json:"image"`
	Manifest   string  `yaml:"manifest" json:"manifest"`
	IsLocal    bool    `yaml:"isLocal" json:"isLocal"`
	StartBlock *uint64 `yaml:"startBlock" json:"startBlock,omitempty"`
	StopBlock  *uint64 `yaml:"stopBlock" json:"stopBlock,omitempty"`
}

// ToAgentInfo transforms the agent config to the agent info.
func (ac AgentConfig) ToAgentInfo() *protocol.AgentInfo {
	return &protocol.AgentInfo{
		Id:        ac.ID,
		Image:     ac.Image,
		ImageHash: ac.ImageHash(),
		Manifest:  ac.Manifest,
	}
}

func (ac AgentConfig) ImageHash() string {
	_, digest := utils.SplitImageRef(ac.Image)
	return digest
}

func (ac AgentConfig) ContainerName() string {
	return ac.containerName("agent")
}

func (ac AgentConfig) AdminContainerName() string {
	return ac.containerName("agent-admin")
}

func (ac AgentConfig) containerName(name string) string {
	_, digest := utils.SplitImageRef(ac.Image)
	parts := []string{ContainerNamePrefix, name, utils.ShortenString(ac.ID, 8)}
	if !ac.IsLocal {
		parts = append(parts, utils.ShortenString(digest, 4))
	}
	return strings.Join(parts, "-")
}

func (ac AgentConfig) GrpcPort() string {
	return AgentGrpcPort
}
