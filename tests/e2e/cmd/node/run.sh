#!/bin/bash

set -xe
set -o pipefail

BIN_NAME="$1"
BIN_NAME=${BIN_NAME:1}
CMD_NAME="$2"

mkdir -p /.forta/coverage
TIMESTAMP=$(date +%s)
COVERAGE_TMP="/.forta/coverage/$CMD_NAME-coverage-$TIMESTAMP.tmp"
touch "$COVERAGE_TMP"

exec "$BIN_NAME" -test.coverprofile="$COVERAGE_TMP" "$CMD_NAME"
