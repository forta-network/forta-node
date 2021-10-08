#!/bin/bash

set -xe

# dynamically look up the secret name
instanceId=$(curl -s http://instance-data/latest/meta-data/instance-id)
region=$(curl -s http://169.254.169.254/latest/dynamic/instance-identity/document | jq -r .region)
envPrefix=$(aws ec2 describe-tags --region $region --filters "Name=resource-id,Values=$instanceId" "Name=key,Values=Environment" |jq -r '.Tags[0].Value')
secretId="${envPrefix}_forta_passphrase"

# get secret from secrets manager
passphrase=$(aws secretsmanager --region $region get-secret-value --secret-id $secretId |jq -r '.SecretString')

# get private key JSON from DynamoDB
privateKeysTable="$envPrefix-forta-node-private-keys"
nodeName=$(aws ec2 describe-tags --region $region --filters "Name=resource-id,Values=$instanceId" "Name=key,Values=Name" | jq -r '.Tags[0].Value')
networkName=$(aws ec2 describe-tags --region $region --filters "Name=resource-id,Values=$instanceId" "Name=key,Values=Network" | jq -r '.Tags[0].Value')
privateKeyItem=$(aws dynamodb get-item --region $region --table $privateKeysTable --key '{"NodeName": { "S": "'$nodeName'" }, "Network": { "S": "'$networkName'"}}' | jq -r .Item)
privateKeyFileName=''
# create and store new one if it doesn't exist
if [ -z "$privateKeyItem" ]; then
	geth account new --password <(echo $passphrase) --keystore "$HOME/.forta/.keys"
	privateKeyFileName=$(ls $HOME/.forta/.keys | head -n 1)
	privateKeyJson=$(cat "$HOME/.forta/.keys/$privateKeyFileName")
	ethereumAddress="0x$(echo $privateKeyJson | jq -r .address)"
	dynamoItemTpl='{NodeName:{S:$name},EthereumAddress:{S:$address},PrivateKeyJson:{S:$privKeyJson},FileName:{S:$keyFileName},Network:{S:$networkName}}'
	privateKeyItem=$(jq -ncM --arg name "$nodeName" --arg address "$ethereumAddress" --arg privKeyJson "$privateKeyJson" --arg keyFileName "$privateKeyFileName" --arg networkName "$networkName" "$dynamoItemTpl")
	aws dynamodb put-item --region $region --table-name $privateKeysTable --item "$privateKeyItem"
fi
privateKeyJson=$(echo $privateKeyItem | jq -r '.PrivateKeyJson.S')
privateKeyFileName=$(echo $privateKeyItem | jq -r '.FileName.S')
# write the private key file to ensure it exists in the right place
mkdir -p "$HOME/.forta/.keys"
echo "$privateKeyJson" > "$HOME/.forta/.keys/$privateKeyFileName"

nohup \
	forta --config "/etc/forta/config.yml" --passphrase $passphrase run > /dev/null 2> /tmp/forta.log < /dev/null &
