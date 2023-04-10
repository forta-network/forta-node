#!/bin/bash

export RUNNER_TRACKING_ID=""

TEST_DIR=$(dirname "${BASH_SOURCE[0]}")
$TEST_DIR/deps-stop.sh

export IPFS_PATH="$TEST_DIR/.ipfs"
export REGISTRY_CONFIGURATION_PATH="$TEST_DIR/disco.config.yml"
export IPFS_URL="http://localhost:5002"
export DISCO_PORT="1970"

ETHEREUM_DIR="$TEST_DIR/.ethereum"
ETHEREUM_PASSWORD_FILE="$TEST_DIR/ethaccounts/password"
ETHEREUM_KEY_FILE="$TEST_DIR/ethaccounts/gethkeyfile"
ETHEREUM_GENESIS_FILE="$TEST_DIR/genesis.json"
ETHEREUM_NODE_ADDRESS="0x1111e291778AE830cfE4e34185e4e560E94047c7"

# ensure that the binaries are installed
set -e
which geth docker ipfs disco
set +e

# spawn mock graphql api
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mock-graphql-api $TEST_DIR/cmd/graphql-api/main.go
./mock-graphql-api &

# ignore error from 'ipfs init' here since it might be failing due to reusing ipfs dir from previous run.
# this is useful for making container-related steps faster in local development.
ipfs init
ipfs daemon --routing none --offline > /dev/null 2>&1 &

disco > /dev/null 2>&1 &

rm -rf "$ETHEREUM_DIR"
geth account import --datadir "$ETHEREUM_DIR" --password "$ETHEREUM_PASSWORD_FILE" "$ETHEREUM_KEY_FILE"
geth init --datadir "$ETHEREUM_DIR" "$ETHEREUM_GENESIS_FILE"
# rpc.gascap=0 means infinite
geth \
	--nodiscover \
	--miner.etherbase $ETHEREUM_NODE_ADDRESS \
	--rpc.allow-unprotected-txs \
	--rpc.gascap 0 \
	--networkid 137 \
	--datadir "$ETHEREUM_DIR" \
	--allow-insecure-unlock \
	--unlock "$ETHEREUM_NODE_ADDRESS" \
	--password "$ETHEREUM_PASSWORD_FILE" \
	--mine \
	--http \
	--http.vhosts '*' \
	--http.port 8545 \
	--http.addr '0.0.0.0' \
	--http.corsdomain '*' \
	--http.api personal,eth,net,web3,txpool,miner \
	> /dev/null 2>&1 &
