#!/bin/bash

set -e
set -o pipefail

REGISTRY="$1"

SCANNER_IMAGE=$(./build_and_push.sh "$REGISTRY" 'scanner')
QUERY_IMAGE=$(./build_and_push.sh "$REGISTRY" 'query')
JSON_RPC_IMAGE=$(./build_and_push.sh "$REGISTRY" 'json-rpc')

./scripts/build.sh "$SCANNER_IMAGE" "$QUERY_IMAGE" "$JSON_RPC_IMAGE" 'remote'
