#!/bin/bash

set -xe
set -o pipefail

REGISTRY="$1"

DIGEST=''
push_and_find_digest() {
	PUSH_OUTPUT=$(docker push "$REGISTRY/$1")
	DIGEST=$(grep -oE '([0-9a-z]{64})' "$PUSH_OUTPUT")
}

docker build -t "$REGISTRY/forta-scanner" -f Dockerfile.scanner .
push_and_find_digest forta-scanner
SCANNER_IMAGE="$REGISTRY/$DIGEST"

docker build -t "$REGISTRY/forta-query" -f Dockerfile.query .
push_and_find_digest forta-query
QUERY_IMAGE="$REGISTRY/$DIGEST"

docker build -t "$REGISTRY/forta-json-rpc" -f Dockerfile.json-rpc .
push_and_find_digest forta-json-rpc
JSON_RPC_IMAGE="$REGISTRY/$DIGEST"

mkdir -p build
go build -o build/forta "-X 'config/config.DockerScannerImageName=$SCANNER_IMAGE' -X 'config/config.DockerQueryContainerName=$QUERY_IMAGE' -X 'config/config.DockerJSONRPCProxyContainerName=$JSON_RPC_IMAGE'"
