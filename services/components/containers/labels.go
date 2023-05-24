package containers

import (
	"github.com/docker/docker/api/types"
	"github.com/forta-network/forta-node/clients/docker"
)

// HasSameLabelValue checks if given container has same label value.
func HasSameLabelValue(container *types.Container, key, value string) bool {
	return container.Labels[key] == value
}

// IsBotContainer checks if given container is a bot container by looking at the label value.
func IsBotContainer(container *types.Container) bool {
	return HasSameLabelValue(container, docker.LabelFortaIsBot, LabelValueFortaIsBot)
}
