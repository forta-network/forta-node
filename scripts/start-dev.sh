#!/bin/bash

source /home/ec2-user/.bash_profile
GITHUB_HASH=$1
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 997179694723.dkr.ecr.us-west-2.amazonaws.com
docker pull 997179694723.dkr.ecr.us-west-2.amazonaws.com/fortify-scanner:$GITHUB_HASH
docker pull 997179694723.dkr.ecr.us-west-2.amazonaws.com/fortify-query:$GITHUB_HASH
docker pull 997179694723.dkr.ecr.us-west-2.amazonaws.com/fortify-json-rpc:$GITHUB_HASH
killall -s SIGINT fortify || true
sleep 3

# this is run over ssh and the path is a directory up
nohup ./fortify/fortify -config fortify/config.yml -passphrase $PASSPHRASE &