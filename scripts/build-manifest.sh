#!/bin/bash

cp $1 $2

ts=$(date --utc +%FT%T.%3NZ)
sed -i "s/%TIMESTAMP%/$ts/g" manifest.json
sed -i "s/%COMMIT_SHA%/$3/g" manifest.json
sed -i "s/%NODE_IMAGE%/$4/g" manifest.json
