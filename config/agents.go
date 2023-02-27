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
	ID          string  `yaml:"id" json:"id"`
	Image       string  `yaml:"image" json:"image"`
	Manifest    string  `yaml:"manifest" json:"manifest"`
	IsLocal     bool    `yaml:"isLocal" json:"isLocal"`
	StartBlock  *uint64 `yaml:"startBlock" json:"startBlock,omitempty"`
	StopBlock   *uint64 `yaml:"stopBlock" json:"stopBlock,omitempty"`
	AlertConfig *protocol.AlertConfig
	ShardConfig *ShardConfig
}

type ShardConfig struct {
	ShardID uint `yaml:"shardId" json:"shardId"`
	Shards  uint `yaml:"shards" json:"shards"`
	Target  uint `yaml:"target" json:"target"`
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
		return fmt.Sprintf("%s-agent-%s", ContainerNamePrefix, utils.ShortenString(ac.ID, 8))
	}
	return fmt.Sprintf(
		"%s-agent-%s-%s", ContainerNamePrefix, utils.ShortenString(ac.ID, 8), utils.ShortenString(digest, 4),
	)
}

func (ac AgentConfig) IsEqual(b AgentConfig) bool {
	sameID := strings.EqualFold(ac.ID, b.ID)
	sameDigest := strings.EqualFold(ac.Image, b.Image)

	return sameID && sameDigest
}

func (ac AgentConfig) GrpcPort() string {
	return AgentGrpcPort
}
