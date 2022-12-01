#!/bin/bash

set -ex
set -o pipefail

TEST_DIR=$(dirname "${BASH_SOURCE[0]}")
SCRIPTS_DIR="$TEST_DIR/../../scripts"

REGISTRY="localhost:1970"
AGENT_IMAGE_SHORT_NAME="forta-e2e-test-agent"
AGENT_IMAGE_FULL_NAME="$REGISTRY/$AGENT_IMAGE_SHORT_NAME"

ALERT_AGENT_IMAGE_SHORT_NAME="forta-e2e-alert-test-agent"
ALERT_AGENT_IMAGE_FULL_NAME="$REGISTRY/$ALERT_AGENT_IMAGE_SHORT_NAME"

# build test agent image
DOCKER_BUILDKIT=1 docker build --network host -t "$AGENT_IMAGE_FULL_NAME" -f \
    "$TEST_DIR/agents/txdetectoragent/Dockerfile" .
docker tag "$AGENT_IMAGE_FULL_NAME" "$AGENT_IMAGE_SHORT_NAME"

# build test alert bot image
DOCKER_BUILDKIT=1 docker build --network host -t "$ALERT_AGENT_IMAGE_FULL_NAME" -f \
    "$TEST_DIR/agents/combinerbot/Dockerfile" .
docker tag "$ALERT_AGENT_IMAGE_FULL_NAME" "$ALERT_AGENT_IMAGE_SHORT_NAME"

AGENT_IMAGE_REF=$("$SCRIPTS_DIR/docker-push.sh" "$REGISTRY" "$AGENT_IMAGE_FULL_NAME")
ALERT_AGENT_IMAGE_REF=$("$SCRIPTS_DIR/docker-push.sh" "$REGISTRY" "$ALERT_AGENT_IMAGE_FULL_NAME")

IMAGE_REFS_DIR="$TEST_DIR/.imagerefs"
mkdir -p "$IMAGE_REFS_DIR"
echo "$AGENT_IMAGE_REF" > "$IMAGE_REFS_DIR/agent"
echo "$ALERT_AGENT_IMAGE_REF" > "$IMAGE_REFS_DIR/combinerbot"