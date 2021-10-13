#!/bin/bash

set -e
set -o pipefail

SCANNER_IMAGE='forta-protocol/forta-scanner:latest'
QUERY_IMAGE='forta-protocol/forta-query:latest'
JSON_RPC_IMAGE='forta-protocol/forta-json-rpc:latest'

docker build -t "$SCANNER_IMAGE" -f Dockerfile.scanner .
docker build -t "$QUERY_IMAGE" -f Dockerfile.query .
docker build -t "$JSON_RPC_IMAGE" -f Dockerfile.json-rpc .

go build -o forta .
