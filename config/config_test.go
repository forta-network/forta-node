package config

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplyContextDefaults(t *testing.T) {
	r := require.New(t)

	cfg := &Config{
		ChainID: 1, // trace enabled
		LocalModeConfig: LocalModeConfig{
			Enable: true, // trace not enforced
		},
		Trace: TraceConfig{
			Enabled: false,
		},
	}

	applyContextDefaults(cfg)

	r.Equal(cfg.Trace.Enabled, false)
	r.Equal(DefaultContainerFortaDirPath, cfg.FortaDir)
	r.Equal(path.Join(cfg.FortaDir, DefaultKeysDirName), cfg.KeyDirPath)
	r.Equal(path.Join(cfg.FortaDir, DefaultCombinerCacheFileName), cfg.CombinerConfig.CombinerCachePath)
}
