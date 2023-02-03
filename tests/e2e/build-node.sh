#!/bin/bash

set -ex
set -o pipefail

# SKIP_CONTAINER_BUILD="$1"
# if [ "$SKIP_CONTAINER_BUILD" == "1" ]; then
# 	exit 0
# fi

TEST_DIR=$(dirname "${BASH_SOURCE[0]}")
ROOT_DIR="$TEST_DIR/../.."
SCRIPTS_DIR="$ROOT_DIR/scripts"

REGISTRY="localhost:1970"
NODE_IMAGE_FULL_NAME="$REGISTRY/forta-node"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o forta-node "$ROOT_DIR/cmd/node/main.go"
DOCKER_BUILDKIT=1 docker build --network=host -t "$NODE_IMAGE_FULL_NAME" -f "$ROOT_DIR/Dockerfile.buildkit.dev.node" .

NODE_IMAGE_REF=$("$SCRIPTS_DIR/docker-push.sh" "$REGISTRY" "$NODE_IMAGE_FULL_NAME")

IMAGE_REFS_DIR="$TEST_DIR/.imagerefs"
mkdir -p "$IMAGE_REFS_DIR"
echo "$NODE_IMAGE_REF" > "$IMAGE_REFS_DIR/node"
