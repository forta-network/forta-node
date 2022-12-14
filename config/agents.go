package config

import (
	"fmt"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
)

const (
	AgentGrpcPort = "50051"
)

type AgentConfig struct {
	ID          string  `yaml:"id" json:"id"`
	Image       string  `yaml:"image" json:"image"`
	Manifest    string  `yaml:"manifest" json:"manifest"`
	IsLocal     bool    `yaml:"isLocal" json:"isLocal"`
	StartBlock  *uint64 `yaml:"startBlock" json:"startBlock,omitempty"`
	StopBlock   *uint64 `yaml:"stopBlock" json:"stopBlock,omitempty"`
	Redundancy  uint    `yaml:"redundancy" json:"redundancy"`
	ShardCount  uint    `yaml:"shardCount" json:"shardCount"`
	ShardId     uint    `yaml:"shardId" json:"shardId"`
	AlertConfig *protocol.AlertConfig
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
	_, digest := utils.SplitImageRef(ac.Image)
	if ac.IsLocal {
		return fmt.Sprintf("%s-agent-%s-%d", ContainerNamePrefix, utils.ShortenString(ac.ID, 8), ac.ShardId)
	}

	return fmt.Sprintf(
		"%s-agent-%s-%s-%d", ContainerNamePrefix, utils.ShortenString(ac.ID, 8),
		utils.ShortenString(digest, 4), ac.ShardId,
	)
}

func (ac AgentConfig) GrpcPort() string {
	return AgentGrpcPort
}
