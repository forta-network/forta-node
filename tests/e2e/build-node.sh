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

# build a node image that creates coverage output
DOCKER_BUILDKIT=1 docker build --network host -t "$NODE_IMAGE_FULL_NAME" -f "$TEST_DIR/cmd/node/Dockerfile" .

NODE_IMAGE_REF=$("$SCRIPTS_DIR/docker-push.sh" "$REGISTRY" "$NODE_IMAGE_FULL_NAME")

IMAGE_REFS_DIR="$TEST_DIR/.imagerefs"
mkdir -p "$IMAGE_REFS_DIR"
echo "$NODE_IMAGE_REF" > "$IMAGE_REFS_DIR/node"