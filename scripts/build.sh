#!/bin/bash

set -e
set -o pipefail

MODULE_NAME=$(grep 'module' go.mod | cut -c8-) # Get the module name from go.mod
IMPORT="$MODULE_NAME/config"
go build -o forta -ldflags="-X '$IMPORT.DockerSupervisorImage=$1' -X '$IMPORT.DockerUpdaterImage=$1' -X '$IMPORT.UseDockerImages=$2' \
 -X '$IMPORT.ReleaseCid=$3' -X '$IMPORT.CommitHash=$4' -X '$IMPORT.Version=$5'" .
