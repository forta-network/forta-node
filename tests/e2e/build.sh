#!/bin/bash

set -ex
set -o pipefail

SKIP_CONTAINER_BUILD="$1"
if [ "$SKIP_CONTAINER_BUILD" == "1" ]; then
	exit 0
fi

REGISTRY="localhost:1970"
NODE_IMAGE_FULL_NAME="$REGISTRY/forta-node"
AGENT_IMAGE_FULL_NAME="$REGISTRY/forta-e2e-test-agent"

# build a node image that creates coverage output
DOCKER_BUILDKIT=1 docker build -t "$NODE_IMAGE_FULL_NAME" -f cmd/node/Dockerfile ./../..
# build test agent image
DOCKER_BUILDKIT=1 docker build -t "$AGENT_IMAGE_FULL_NAME" -f agents/txdetectoragent/Dockerfile ./../..

NODE_IMAGE_REF=$(../../scripts/docker-push.sh "$REGISTRY" "$NODE_IMAGE_FULL_NAME")
AGENT_IMAGE_REF=$(../../scripts/docker-push.sh "$REGISTRY" "$AGENT_IMAGE_FULL_NAME")

mkdir -p .imagerefs

echo "$NODE_IMAGE_REF" > .imagerefs/node
echo "$AGENT_IMAGE_REF" > .imagerefs/agent
