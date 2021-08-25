#!/bin/bash

set -e
set -o pipefail

mkdir -p build
MODULE_NAME=$(grep 'module' go.mod | cut -c8-) # Get the module name from go.mod
IMPORT="$MODULE_NAME/config"
go build -o forta -ldflags="-X '$IMPORT.DockerScannerContainerImage=$1' -X '$IMPORT.DockerQueryContainerImage=$2' -X '$IMPORT.DockerJSONRPCProxyContainerImage=$3' -X '$IMPORT.UseDockerImages=$4'"
