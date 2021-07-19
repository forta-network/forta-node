#!/bin/bash

# this looks up the alchemy credentials and injects them into the config yaml
instanceId=$(curl -s http://instance-data/latest/meta-data/instance-id)
region=$(curl -s http://169.254.169.254/latest/dynamic/instance-identity/document | jq -r .region)
envPrefix=$(aws ec2 describe-tags --region $region --filters "Name=resource-id,Values=$instanceId" "Name=key,Values=Environment" |jq -r '.Tags[0].Value')
secretId="${envPrefix}_alchemy_api_url"

apiUrlUnsafe=$(aws secretsmanager --region $region get-secret-value --secret-id $secretId |jq -r '.SecretString')
apiUrl=$(printf '%s\n' "$apiUrlUnsafe" | sed -e 's/[]\/$*.^[]/\\&/g');

sed -i "s/ALCHEMY_URL/$apiUrl/g" /etc/fortify/config-fortify.yml
