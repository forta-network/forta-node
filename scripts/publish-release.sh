#!/bin/bash

hash=$1
key=$2
url=$3

curl -f -X POST -H "x-forta-key: $key" -H "content-type: application/json" -d "{\"hash\": \"$hash\"}" "$3"
