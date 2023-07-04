#!/bin/bash
echo "***FORTA START***"

echo "Getting latest block..." && \
FORK_BLOCK=`cast block --rpc-url $RPC_URL | grep "number" | grep -Eo '[0-9]{8}'` && \

echo "Starting at block $FORK_BLOCK" && \
START_BLOCK=$FORK_BLOCK docker compose -f $FORTA/docker-compose.yml up --build

exit 0