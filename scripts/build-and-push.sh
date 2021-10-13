#!/bin/bash

set -e

REGISTRY="$1"
IMAGE_NAME="$2"

push_and_find_digest() {
	PUSH_OUTPUT=$(docker push "$REGISTRY/$1")
	DIGEST=$(echo "$PUSH_OUTPUT" | grep -oE '([0-9a-z]{64})')
}

docker build -t "$REGISTRY/forta-$IMAGE_NAME" -f "Dockerfile.$IMAGE_NAME" .
push_and_find_digest "$IMAGE_NAME"
echo "$REGISTRY/$DIGEST"
