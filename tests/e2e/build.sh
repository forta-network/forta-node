#!/bin/bash

set -ex
set -o pipefail

# SKIP_CONTAINER_BUILD="$1"
# if [ "$SKIP_CONTAINER_BUILD" == "1" ]; then
# 	exit 0
# fi

TEST_DIR=$(dirname "${BASH_SOURCE[0]}")
SCRIPTS_DIR="$TEST_DIR/../../scripts"

REGISTRY="localhost:1970"
NODE_IMAGE_FULL_NAME="$REGISTRY/forta-node"

# build bot images
bash "$(dirname "$0")"/build-bots.sh

# build a node image that creates coverage output
bash "$(dirname "$0")"/build-node.sh

NODE_IMAGE_REF=$("$SCRIPTS_DIR/docker-push.sh" "$REGISTRY" "$NODE_IMAGE_FULL_NAME")

IMAGE_REFS_DIR="$TEST_DIR/.imagerefs"
mkdir -p "$IMAGE_REFS_DIR"
echo "$NODE_IMAGE_REF" > "$IMAGE_REFS_DIR/node"

# build the test cli/runner binary
MODULE_NAME=$(grep 'module' "$TEST_DIR/../../go.mod" | cut -c8-) # Get the module name from go.mod
IMPORT="$MODULE_NAME/config"
GO_PACKAGES=$(go list ./... | grep -v tests | tr "\n" ",")
GO_PACKAGES=${GO_PACKAGES%?} # cut trailing comma

go test -c -o forta-test \
	-ldflags="-X '$IMPORT.DockerSupervisorImage=$NODE_IMAGE_REF' -X '$IMPORT.DockerUpdaterImage=$NODE_IMAGE_REF' -X '$IMPORT.UseDockerImages=remote' -X '$IMPORT.Version=0.0.1-test'" \
	"$TEST_DIR/cmd/cli"
mv -f forta-test "$TEST_DIR/"
