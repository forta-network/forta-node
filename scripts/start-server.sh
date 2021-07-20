#!/bin/bash

# dynamically look up the secret name
instanceId=$(curl -s http://instance-data/latest/meta-data/instance-id)
region=$(curl -s http://169.254.169.254/latest/dynamic/instance-identity/document | jq -r .region)
envPrefix=$(aws ec2 describe-tags --region $region --filters "Name=resource-id,Values=$instanceId" "Name=key,Values=Environment" |jq -r '.Tags[0].Value')
secretId="${envPrefix}_fortify_passphrase"

# get secret from secrets manager
passphrase=$(aws secretsmanager --region $region get-secret-value --secret-id $secretId |jq -r '.SecretString')
nohup fortify -config "/etc/fortify/config-fortify.yml" -passphrase $passphrase > /dev/null 2> /tmp/forta.log < /dev/null &