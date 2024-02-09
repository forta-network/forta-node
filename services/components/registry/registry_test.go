package registry

import (
	"errors"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/config"
	mock_store "github.com/forta-network/forta-node/store/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestLoadAssignedBots(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	regStore := mock_store.NewMockRegistryStore(ctrl)
	botReg := &botRegistry{
		scannerAddress: common.HexToAddress(utils.ZeroAddress),
		registryStore:  regStore,
		mu:             &sync.RWMutex{},
	}

	cfgs := []config.AgentConfig{{}}
	regStore.EXPECT().GetAgentsIfChanged(utils.ZeroAddress).Return(cfgs, true, nil)
	retCfgs, err := botReg.LoadAssignedBots()
	r.NoError(err)
	r.Equal(cfgs, retCfgs)
	r.Equal(cfgs, botReg.botConfigs)

	changedCfg := []config.AgentConfig{{}, {}}
	regStore.EXPECT().GetAgentsIfChanged(utils.ZeroAddress).Return(changedCfg, false, nil)
	retCfgs, err = botReg.LoadAssignedBots()
	r.NoError(err)
	r.Equal(cfgs, retCfgs)
	r.Equal(cfgs, botReg.botConfigs)

	regStore.EXPECT().GetAgentsIfChanged(utils.ZeroAddress).Return(nil, false, errors.New("some error"))
	retCfgs, err = botReg.LoadAssignedBots()
	r.Error(err)
	r.Nil(retCfgs)
}

func TestBotRegistry_LoadHeartbeatBot(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	regStore := mock_store.NewMockRegistryStore(ctrl)
	botReg := &botRegistry{
		scannerAddress: common.HexToAddress(utils.ZeroAddress),
		registryStore:  regStore,
	}
	cfg := &config.AgentConfig{ID: config.HeartbeatBotID}
	regStore.EXPECT().FindAgentGlobally(config.HeartbeatBotID).Return(cfg, nil)

	res, err := botReg.LoadHeartbeatBot()
	r.NoError(err)
	r.Equal(config.HeartbeatBotID, res.ID)
}
