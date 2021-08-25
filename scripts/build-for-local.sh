#!/bin/bash

set -e
set -o pipefail

SCANNER_IMAGE='forta-network/forta-scanner:latest'
QUERY_IMAGE='forta-network/forta-query:latest'
JSON_RPC_IMAGE='forta-network/forta-json-rpc:latest'

docker build -t "$SCANNER_IMAGE" -f Dockerfile.scanner .
docker build -t "$QUERY_IMAGE" -f Dockerfile.query .
docker build -t "$JSON_RPC_IMAGE" -f Dockerfile.json-rpc .

./scripts/build.sh "$SCANNER_IMAGE" "$QUERY_IMAGE" "$JSON_RPC_IMAGE" 'local'
