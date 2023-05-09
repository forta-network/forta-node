package config

import (
	"fmt"
	"strings"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
)

const (
	AgentGrpcPort = "50051"
)

type AgentConfig struct {
	ID           string  `yaml:"id" json:"id"`
	Image        string  `yaml:"image" json:"image"`
	Manifest     string  `yaml:"manifest" json:"manifest"`
	IsLocal      bool    `yaml:"isLocal" json:"isLocal"`
	IsStandalone bool    `yaml:"isStandalone" json:"isStandalone"`
	StartBlock   *uint64 `yaml:"startBlock" json:"startBlock,omitempty"`
	StopBlock    *uint64 `yaml:"stopBlock" json:"stopBlock,omitempty"`
	Owner        string  `yaml:"owner" json:"owner"`

	ChainID     int
	AlertConfig *protocol.AlertConfig
	ShardConfig *ShardConfig
}

type ShardConfig struct {
	ShardID uint `yaml:"shardId" json:"shardId"`
	Shards  uint `yaml:"shards" json:"shards"`
	Target  uint `yaml:"target" json:"target"`
}

func (ac AgentConfig) Equal(b AgentConfig) bool {
	sameID := strings.EqualFold(ac.ID, b.ID)
	sameManifest := strings.EqualFold(ac.Manifest, b.Manifest)
	if !sameID || !sameManifest {
		return false
	}

	noSharding := ac.ShardConfig == nil && b.ShardConfig == nil
	if noSharding {
		return true
	}
	sameShardID := ac.ShardConfig.ShardID == b.ShardConfig.ShardID
	sameShardCount := ac.ShardConfig.Shards == b.ShardConfig.Shards

	return sameShardID && sameShardCount
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
	if ac.IsStandalone {
		// the container is already running - don't mess with the name
		return ac.ID
	}
	if ac.IsLocal {
		return fmt.Sprintf("%s-agent-%s", ContainerNamePrefix, utils.ShortenString(ac.ID, 8))
	}
	_, digest := utils.SplitImageRef(ac.Image)
	return fmt.Sprintf(
		"%s-agent-%s-%s", ContainerNamePrefix, utils.ShortenString(ac.ID, 8), utils.ShortenString(digest, 4),
	)
}

func (ac AgentConfig) GrpcPort() string {
	return AgentGrpcPort
}
