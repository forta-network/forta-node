package config

import (
	"fmt"
	"path"
)

const ContainerNamePrefix = "forta"

// Docker container names
var (
	DockerScannerNodeImage = "forta-protocol/forta-node:latest"
	UseDockerImages        = "local"

	DockerSupervisorContainerName   = fmt.Sprintf("%s-supervisor", ContainerNamePrefix)
	DockerNatsContainerName         = fmt.Sprintf("%s-nats", ContainerNamePrefix)
	DockerScannerContainerName      = fmt.Sprintf("%s-scanner", ContainerNamePrefix)
	DockerJSONRPCProxyContainerName = fmt.Sprintf("%s-json-rpc", ContainerNamePrefix)
	DockerPublisherContainerName    = fmt.Sprintf("%s-publisher", ContainerNamePrefix)

	DockerNetworkName = DockerScannerContainerName

	DefaultContainerFortaDirPath        = "/.forta"
	DefaultContainerKeyDirPath          = path.Join(DefaultContainerFortaDirPath, DefaultKeysDirName)
	DefaultContainerLocalAgentsFilePath = path.Join(DefaultContainerFortaDirPath, DefaultLocalAgentsFileName)
	DefaultContainerConfigPath          = "/forta_config.yml"
)
