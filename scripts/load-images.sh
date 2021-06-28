#!/bin/bash

# ECR_REGISTRY and GITHUB_HASH are replaced during code deploy
aws ecr get-login-password --region $region | docker login --username AWS --password-stdin "ECR_REGISTRY"
docker pull "ECR_REGISTRY/fortify-scanner:GITHUB_HASH"
docker pull "ECR_REGISTRY/fortify-query:GITHUB_HASH"
docker pull "ECR_REGISTRY/fortify-json-rpc:GITHUB_HASH"


instanceId=$(curl -s http://instance-data/latest/meta-data/instance-id)
region=$(curl -s http://169.254.169.254/latest/dynamic/instance-identity/document | jq -r .region)
envPrefix=$(aws ec2 describe-tags --region $region --filters "Name=resource-id,Values=$instanceId" "Name=key,Values=Environment" |jq -r '.Tags[0].Value')

# pull agents by parsing config file
cat "/etc/fortify/config-fortify-${envPrefix}.yml" | grep image | sed -e 's/.* //g'|xargs -I {} docker pull {}