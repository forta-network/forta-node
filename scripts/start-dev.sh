#!/bin/bash

GITHUB_HASH=$1
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 997179694723.dkr.ecr.us-west-2.amazonaws.com
docker pull 997179694723.dkr.ecr.us-west-2.amazonaws.com/fortify-scanner:$GITHUB_HASH
docker pull 997179694723.dkr.ecr.us-west-2.amazonaws.com/fortify-query:$GITHUB_HASH
docker pull 997179694723.dkr.ecr.us-west-2.amazonaws.com/fortify-json-rpc:$GITHUB_HASH
killall -s SIGINT fortify
sleep 3
nohup ./fortify -passphrase $PASSPHRASE &
