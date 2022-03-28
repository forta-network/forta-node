#!/bin/bash

BIN_NAME="$1"
CMD_NAME="$2"

"${BIN_NAME:1}" "$CMD_NAME" -test.coverprofile=coverage.tmp && tail -n +2 coverage.tmp >> /.forta/coverage.txt
