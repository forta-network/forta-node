package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
)

const (
	AgentGrpcPort  = "50051"
	HeartbeatBotID = "0x37ce3df7facea4bd95619655347ffd66f50909d100fa74ae042f122a8b024c29"
)

type AgentConfig struct {
	ID              string  `yaml:"id" json:"id"`
	Image           string  `yaml:"image" json:"image"`
	Manifest        string  `yaml:"manifest" json:"manifest"`
	IsLocal         bool    `yaml:"isLocal" json:"isLocal"`
	IsStandalone    bool    `yaml:"isStandalone" json:"isStandalone"`
	StartBlock      *uint64 `yaml:"startBlock" json:"startBlock,omitempty"`
	StopBlock       *uint64 `yaml:"stopBlock" json:"stopBlock,omitempty"`
	Owner           string  `yaml:"owner" json:"owner"`
	ProtocolVersion int     `yaml:"protocolVersion" json:"protocolVersion"`

	ChainID     int
	ShardConfig *ShardConfig
}

type ShardConfig struct {
	ShardID uint `yaml:"shardId" json:"shardId"`
	Shards  uint `yaml:"shards" json:"shards"`
	Target  uint `yaml:"target" json:"target"`
}

func (ac AgentConfig) ShardID() int32 {
	if !ac.IsSharded() {
		// default uint value is 0, so we cannot tell between an unset shardID and actual shard 0
		// therefore, we return -1 here
		return -1
	}
	return int32(ac.ShardConfig.ShardID)
}

func (ac AgentConfig) ShardDetails() string {
	if !ac.IsSharded() {
		return ""
	}
	return fmt.Sprintf("shard=%d", ac.ShardConfig.ShardID)
}

func (ac AgentConfig) Equal(b AgentConfig) bool {
	sameID := strings.EqualFold(ac.ID, b.ID)
	sameManifest := strings.EqualFold(ac.Manifest, b.Manifest)
	if !sameID || !sameManifest {
		return false
	}

	// if both don't have sharding config, then they are the same
	if ac.ShardConfig == nil && b.ShardConfig == nil {
		return true
	}

	// if one of them does not have shard config, then they are different
	if ac.ShardConfig == nil && b.ShardConfig != nil {
		return false
	}

	if ac.ShardConfig != nil && b.ShardConfig == nil {
		return false
	}

	// if both have shard config, then configs should match

	sameShardID := ac.ShardConfig.ShardID == b.ShardConfig.ShardID
	sameShardCount := ac.ShardConfig.Shards == b.ShardConfig.Shards

	return sameShardID && sameShardCount
}

// IsSharded tells if this is a sharded bot.
func (ac *AgentConfig) IsSharded() bool {
	return ac.ShardConfig != nil && ac.ShardConfig.Shards > 1
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

	parts := []string{ContainerNamePrefix, "agent", utils.ShortenString(ac.ID, 8), utils.ShortenString(digest, 4)}
	if ac.IsSharded() {
		parts = append(parts, strconv.Itoa(int(ac.ShardConfig.ShardID))) // append the shard id at the end
	}
	return strings.Join(parts, "-")
}

func (ac AgentConfig) GrpcPort() string {
	return AgentGrpcPort
}
