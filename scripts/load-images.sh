#!/bin/bash

source /home/fortify/.bash_profile

region=$(curl -s http://169.254.169.254/latest/dynamic/instance-identity/document | jq -r .region)
accountId=$(curl -s http://169.254.169.254/latest/dynamic/instance-identity/document | jq -r .accountId)

aws ecr get-login-password --region $region | docker login --username AWS --password-stdin "${accountId}.dkr.ecr.${region}.amazonaws.com"
docker pull "${accountId}.dkr.ecr.${region}.amazonaws.com/fortify-scanner:GITHUB_HASH"
docker pull "${accountId}.dkr.ecr.${region}.amazonaws.com/fortify-query:GITHUB_HASH"
docker pull "${accountId}.dkr.ecr.${region}.amazonaws.com/fortify-json-rpc:GITHUB_HASH"
