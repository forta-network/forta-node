package registry

import (
	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/contracts"
	"OpenZeppelin/fortify-node/domain"
	"OpenZeppelin/fortify-node/services/registry/regtypes"
	"OpenZeppelin/fortify-node/utils"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
)

func (rs *RegistryService) detectAgentEvents(evt *domain.TransactionEvent) error {
	update, agentID, ref, err := rs.detectAgentEvent(evt)
	if err != nil {
		return err
	}
	return rs.sendAgentUpdate(update, agentID, ref)
}

func (rs *RegistryService) detectAgentEvent(evt *domain.TransactionEvent) (update *agentUpdate, agentID [32]byte, ref string, err error) {
	for _, logEntry := range evt.Receipt.Logs {
		ethLog := transformLog(&logEntry)

		var addedEvent *contracts.AgentRegistryAgentAdded
		addedEvent, err = rs.logUnpacker.UnpackAgentRegistryAgentAdded(ethLog)
		if err == nil {
			if (common.Hash)(addedEvent.PoolId).String() != rs.poolID.String() {
				continue
			}
			return &agentUpdate{IsCreation: true}, addedEvent.AgentId, addedEvent.Ref, nil
		}

		var updatedEvent *contracts.AgentRegistryAgentUpdated
		updatedEvent, err = rs.logUnpacker.UnpackAgentRegistryAgentUpdated(ethLog)
		if err == nil {
			if (common.Hash)(updatedEvent.PoolId).String() != rs.poolID.String() {
				continue
			}
			return &agentUpdate{IsUpdate: true}, updatedEvent.AgentId, updatedEvent.Ref, nil
		}

		var removedEvent *contracts.AgentRegistryAgentRemoved
		removedEvent, err = rs.logUnpacker.UnpackAgentRegistryAgentRemoved(ethLog)
		if err == nil {
			if (common.Hash)(removedEvent.PoolId).String() != rs.poolID.String() {
				continue
			}
			return &agentUpdate{IsRemoval: true}, removedEvent.AgentId, "", nil
		}
	}
	update = nil
	err = nil
	return
}

func (rs *RegistryService) sendAgentUpdate(update *agentUpdate, agentID [32]byte, ref string) error {
	agentCfg, err := rs.makeAgentConfig(agentID, ref)
	if err != nil {
		return err
	}
	update.Config = agentCfg
	log.Infof("sending agent update: %+v", update)
	rs.agentUpdates <- update
	return nil
}

func (rs *RegistryService) makeAgentConfig(agentID [32]byte, ref string) (agentCfg config.AgentConfig, err error) {
	agentCfg.ID = (common.Hash)(agentID).String()
	if len(ref) == 0 {
		return
	}

	var (
		agentData *regtypes.AgentFile
	)
	for i := 0; i < 10; i++ {
		agentData, err = rs.ipfsClient.GetAgentFile(ref)
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("failed to load the agent file using ipfs ref: %v", err)
		return
	}

	var ok bool
	agentCfg.Image, ok = utils.ValidateImageRef(rs.cfg.Registry.ContainerRegistry, agentData.Manifest.ImageReference)
	if !ok {
		log.Warnf("invalid agent reference - skipping: %s", agentCfg.Image)
	}

	return
}

func transformLog(log *domain.LogEntry) *types.Log {
	transformed := &types.Log{
		Address:     common.HexToAddress(utils.String(log.Address)),
		Data:        common.FromHex(utils.String(log.Data)),
		BlockHash:   common.HexToHash(utils.String(log.BlockHash)),
		BlockNumber: hexutil.MustDecodeBig(utils.String(log.BlockNumber)).Uint64(),
		TxHash:      common.HexToHash(utils.String(log.TransactionHash)),
		TxIndex:     uint(hexutil.MustDecodeBig(utils.String(log.TransactionIndex)).Uint64()),
		Index:       uint(hexutil.MustDecodeBig(utils.String(log.LogIndex)).Uint64()),
	}
	for _, topic := range log.Topics {
		transformed.Topics = append(transformed.Topics, common.HexToHash(*topic))
	}
	return transformed
}

func (rs *RegistryService) listenToAgentUpdates() {
	for update := range rs.agentUpdates {
		rs.agentUpdatesWg.Wait()
		rs.handleAgentUpdate(update)
		rs.msgClient.Publish(messaging.SubjectAgentsVersionsLatest, rs.agentsConfigs)
	}
	close(rs.done)
}

func (rs *RegistryService) handleAgentUpdate(update *agentUpdate) {
	switch {
	case update.IsCreation:
		// Skip if we already have this agent.
		for _, agent := range rs.agentsConfigs {
			if agent.ID == update.Config.ID {
				return
			}
		}
		rs.agentsConfigs = append(rs.agentsConfigs, &update.Config)

	case update.IsUpdate:
		for _, agent := range rs.agentsConfigs {
			if agent.ID == update.Config.ID {
				agent.Image = update.Config.Image
				return
			}
		}

	case update.IsRemoval:
		var newAgents []*config.AgentConfig
		for _, agent := range rs.agentsConfigs {
			if agent.ID != update.Config.ID {
				newAgents = append(newAgents, agent)
			}
		}
		rs.agentsConfigs = newAgents

	default:
		log.Panicf("tried to handle unknown agent update")
	}
}
