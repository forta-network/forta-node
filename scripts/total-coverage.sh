#!/bin/bash

set -xe
set -o pipefail

COVERAGE_PREFIX="$1"
DEFAULT_NAME="coverage"

if [ ! -z "$COVERAGE_PREFIX" ]; then
	DEFAULT_NAME="$COVERAGE_PREFIX-$DEFAULT_NAME"
fi

ALL_COVERAGE="$DEFAULT_NAME.txt"

echo 'mode: atomic' > $ALL_COVERAGE

for f in coverage/*
do
	tail -n +2 "$f" >> $ALL_COVERAGE
done

go tool cover -func="$ALL_COVERAGE" > "$DEFAULT_NAME.out"
