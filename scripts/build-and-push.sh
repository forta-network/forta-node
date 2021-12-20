#!/bin/bash

set -xe

REGISTRY="$1"
IMAGE_NAME="$2"
FULL_IMAGE_NAME="$REGISTRY/forta-$IMAGE_NAME"

docker build -t "$FULL_IMAGE_NAME" -f "Dockerfile.$IMAGE_NAME" . > /dev/null
PUSH_OUTPUT=$(docker push "$FULL_IMAGE_NAME")
DIGEST=$(echo "$PUSH_OUTPUT" | grep -oE '([0-9a-z]{64})')

FOUND_LINE=$(docker pull -a "$REGISTRY/$DIGEST" | grep bafybei)
if [ -z "$FOUND_LINE" ]; then
	echo "failed to find the IPFS reference of the disco image - exiting"
	exit 1
fi;
IPFS_REF=${FOUND_LINE::59}
echo "$REGISTRY/$IPFS_REF@sha256:$DIGEST"
