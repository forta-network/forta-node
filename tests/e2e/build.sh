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
AGENT_IMAGE_FULL_NAME="$REGISTRY/forta-e2e-test-agent"

# build a node image that creates coverage output
DOCKER_BUILDKIT=1 docker build -t "$NODE_IMAGE_FULL_NAME" -f "$TEST_DIR/cmd/node/Dockerfile" .
# build test agent image
DOCKER_BUILDKIT=1 docker build -t "$AGENT_IMAGE_FULL_NAME" -f "$TEST_DIR/agents/txdetectoragent/Dockerfile" .

NODE_IMAGE_REF=$("$SCRIPTS_DIR/docker-push.sh" "$REGISTRY" "$NODE_IMAGE_FULL_NAME")
AGENT_IMAGE_REF=$("$SCRIPTS_DIR/docker-push.sh" "$REGISTRY" "$AGENT_IMAGE_FULL_NAME")

IMAGE_REFS_DIR="$TEST_DIR/.imagerefs"
mkdir -p "$IMAGE_REFS_DIR"
echo "$NODE_IMAGE_REF" > "$IMAGE_REFS_DIR/node"
echo "$AGENT_IMAGE_REF" > "$IMAGE_REFS_DIR/agent"

# build the test cli/runner binary
MODULE_NAME=$(grep 'module' "$TEST_DIR/../../go.mod" | cut -c8-) # Get the module name from go.mod
IMPORT="$MODULE_NAME/config"
GO_PACKAGES=$(go list ./... | grep -v tests | tr "\n" ",")
GO_PACKAGES=${GO_PACKAGES%?} # cut trailing comma

go test -c -o forta-test -race -covermode=atomic -coverpkg \
	"$GO_PACKAGES" \
	-ldflags="-X '$IMPORT.DockerSupervisorImage=$NODE_IMAGE_REF' -X '$IMPORT.DockerUpdaterImage=$NODE_IMAGE_REF' -X '$IMPORT.UseDockerImages=remote' -X '$IMPORT.Version=0.0.1-test'" \
	"$TEST_DIR/cmd/cli"
mv -f forta-test "$TEST_DIR/"
