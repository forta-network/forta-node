#!/bin/bash

set -e

REGISTRY="$1"
IMAGE_NAME="$2"
FULL_IMAGE_NAME="$REGISTRY/forta-$IMAGE_NAME"

docker build -t "$FULL_IMAGE_NAME" -f "Dockerfile.$IMAGE_NAME" . > /dev/null
PUSH_OUTPUT=$(docker push "$FULL_IMAGE_NAME")
DIGEST=$(echo "$PUSH_OUTPUT" | grep -oE '([0-9a-z]{64})')
echo "$REGISTRY/$DIGEST"
