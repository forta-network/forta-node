#!/bin/bash

# this looks up the alchemy credentials and injects them into the config yaml
instanceId=$(curl -s http://instance-data/latest/meta-data/instance-id)
region=$(curl -s http://169.254.169.254/latest/dynamic/instance-identity/document | jq -r .region)
envPrefix=$(aws ec2 describe-tags --region $region --filters "Name=resource-id,Values=$instanceId" "Name=key,Values=Environment" |jq -r '.Tags[0].Value')

secretId="${envPrefix}_alchemy_api_url"
apiUrlUnsafe=$(aws secretsmanager --region $region get-secret-value --secret-id $secretId |jq -r '.SecretString')
apiUrl=$(printf '%s\n' "$apiUrlUnsafe" | sed -e 's/[]\/$*.^[]/\\&/g');

registryApiUrlId="${envPrefix}_agent_registry_api_url"
registryApiUrlUnsafe=$(aws secretsmanager --region $region get-secret-value --secret-id $registryApiUrlId |jq -r '.SecretString')
registryApiUrl=$(printf '%s\n' "$registryApiUrlUnsafe" | sed -e 's/[]\/$*.^[]/\\&/g');

registryWssUrlId="${envPrefix}_agent_registry_wss_url"
registryWssUrlUnsafe=$(aws secretsmanager --region $region get-secret-value --secret-id $registryWssUrlId |jq -r '.SecretString')
registryWssUrl=$(printf '%s\n' "$registryWssUrlUnsafe" | sed -e 's/[]\/$*.^[]/\\&/g');

sed -i "s/ALCHEMY_URL/$apiUrl/g" /etc/fortify/config-fortify.yml
sed -i "s/REGISTRY_API_URL/$registryApiUrl/g" /etc/fortify/config-fortify.yml
sed -i "s/REGISTRY_WSS_URL/$registryWssUrl/g" /etc/fortify/config-fortify.yml