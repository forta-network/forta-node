#!/bin/bash

set -e
set -o pipefail

REGISTRY="$1"

NODE_IMAGE=$(./build_and_push.sh "$REGISTRY" 'node')

./scripts/build.sh "$NODE_IMAGE" 'remote'
