#!/bin/bash

set -xe

REGISTRY="$1"
IMAGE_NAME="$2"
COMMIT_SHA="$3"
FULL_IMAGE_NAME="$REGISTRY/forta-$IMAGE_NAME-$COMMIT_SHA"

docker build -t "$FULL_IMAGE_NAME" -f "Dockerfile.$IMAGE_NAME" . > /dev/null
./scripts/docker-push.sh "$REGISTRY" "$FULL_IMAGE_NAME"
