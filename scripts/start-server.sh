#!/bin/bash

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
privateKeyItem=$(aws dynamodb get-item --region $region --table $privateKeysTable --key '{"NodeName": { "S": '"$nodeName"' }}' | jq -r .Item)
privateKeyFileName=''
# create and store new one if it doesn't exist
if [ -z "$privateKeyItem" ]; then
	wget -O geth.tar.gz https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.10.9-eae3b194.tar.gz
	tar -xvf geth.tar.gz
	./geth/geth account new --password <(echo $passphrase) --keystore "$HOME/.forta/.keys"
	privateKeyFileName=$(ls ~/.forta/.keys | head -n 1)
	privateKeyJson=$(cat "$HOME/.forta/.keys/$privateKeyFileName")
	ethereumAddress="0x$(echo $privateKeyJson | jq -r .address)"
	privateKeyJson=$(echo $privateKeyJson | jq -RM) # escape so we can put it into another JSON
	dynamoItemTpl='{NodeName:{S:$name},EthereumAddress:{S:$address},PrivateKeyJson:{S:$privKeyJson},FileName:{S:$keyFileName}}'
	privateKeyItem=$(jq -ncM --arg name "$nodeName" --arg address "$ethereumAddress" --arg privKeyJson "$privateKeyJson" --arg keyFileName "$privateKeyFileName" "$dynamoItemTpl")
    aws dynamodb put-item --region $region --table-name $privateKeysTable --item "$privateKeyItem"
fi
privateKeyJson=$(echo $privateKeyItem | jq -r '.PrivateKeyJson.S')
privateKeyFileName=$(echo $privateKeyItem | jq -r '.FileName.S')
# write the private key file to ensure it exists in the right place
cat "$privateKeyJson" > "$HOME/.forta/.keys/$privateKeyFileName"

# get config file name
configFileName=$(aws ec2 describe-tags --region $region --filters "Name=resource-id,Values=$instanceId" "Name=key,Values=FortaConfig" | jq -r '.Tags[0].Value')

nohup \
	forta --config "/etc/forta/configs/$configFileName" --passphrase $passphrase run > /dev/null 2> /tmp/forta.log < /dev/null &
