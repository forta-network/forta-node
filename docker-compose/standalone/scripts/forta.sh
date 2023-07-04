#! /bin/bash
echo "***FORTA***"

echo "Getting latest block..." && \
FORK_BLOCK=`cast block --rpc-url $RPC_URL | grep "number" | grep -Eo '[0-9]{8}'` && \

echo "Starting at block $FORK_BLOCK" && \
echo "RPC URL $RPC_URL" && \
RPC_URL=$RPC_URL FORK_BLOCK=$FORK_BLOCK docker compose -f $FORTA/docker-compose.yml up --remove-orphans --abort-on-container-exit --build


exit 0