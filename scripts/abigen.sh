#!/bin/sh

CONTRACT_NAME_SNAKE_CASE="$1"
CONTRACT_NAME_TITLE_CASE=$(echo "$CONTRACT_NAME_SNAKE_CASE" | sed -r 's/(^|_)(\w)/\U\2/g')
FLAGS="$2"

# filter out error types because they are reported to cause a runtime problem
jq -c '[ .[] | select( .type != "error") ]' \
	< "./contracts/contract_$CONTRACT_NAME_SNAKE_CASE/$CONTRACT_NAME_SNAKE_CASE.json" \
	| abigen --out "./contracts/contract_$CONTRACT_NAME_SNAKE_CASE/$CONTRACT_NAME_SNAKE_CASE.go" \
		--pkg "contract_$CONTRACT_NAME_SNAKE_CASE" \
		--type "$CONTRACT_NAME_TITLE_CASE" \
		$FLAGS \
		--abi -
