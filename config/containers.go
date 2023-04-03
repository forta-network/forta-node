package config

import (
	"fmt"
	"path"
)

const ContainerNamePrefix = "forta"

// Docker container names
var (
	DockerSupervisorImage = "forta-network/forta-node:latest"
	DockerUpdaterImage    = "forta-network/forta-node:latest"
	UseDockerImages       = "local"

	DockerSupervisorManagedContainers = 8
	DockerUpdaterContainerName        = fmt.Sprintf("%s-updater", ContainerNamePrefix)
	DockerSupervisorContainerName     = fmt.Sprintf("%s-supervisor", ContainerNamePrefix)
	DockerNatsContainerName           = fmt.Sprintf("%s-nats", ContainerNamePrefix)
	DockerIpfsContainerName           = fmt.Sprintf("%s-ipfs", ContainerNamePrefix)
	DockerScannerContainerName        = fmt.Sprintf("%s-scanner", ContainerNamePrefix)
	DockerInspectorContainerName      = fmt.Sprintf("%s-inspector", ContainerNamePrefix)
	DockerJSONRPCProxyContainerName   = fmt.Sprintf("%s-json-rpc", ContainerNamePrefix)
	DockerPublicAPIProxyContainerName = fmt.Sprintf("%s-public-api", ContainerNamePrefix)
	DockerJWTProviderContainerName    = fmt.Sprintf("%s-jwt-provider", ContainerNamePrefix)
	DockerStorageContainerName        = fmt.Sprintf("%s-storage", ContainerNamePrefix)

	DockerNetworkName = DockerScannerContainerName

	DefaultContainerFortaDirPath      = "/.forta"
	DefaultContainerConfigPath        = path.Join(DefaultContainerFortaDirPath, DefaultConfigFileName)
	DefaultContainerWrappedConfigPath = path.Join(DefaultContainerFortaDirPath, DefaultWrappedConfigFileName)
	DefaultContainerKeyDirPath        = path.Join(DefaultContainerFortaDirPath, DefaultKeysDirName)
)
