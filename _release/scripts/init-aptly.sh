#!/bin/sh

set -x

AWS_REGION="$1"
AWS_BUCKET_NAME="$2"
AWS_ACCESS_KEY="$3"
AWS_SECRET_KEY="$4"
MIRROR_URL="$5"

cd "$HOME" || exit 1

aptly mirror create fortamirror "$MIRROR_URL" stable # we refer to where we publish as a mirror here, to load light packages
aptly mirror update fortamirror
aptly repo create forta
aptly repo import fortamirror forta forta

set -e

S3_CONFIG="{
	\"releaseBucket\": {
		\"region\": \"$AWS_REGION\",
		\"bucket\": \"$AWS_BUCKET_NAME\",
		\"awsAccessKeyID\": \"$AWS_ACCESS_KEY\",
		\"awsSecretAccessKey\": \"$AWS_SECRET_KEY\"
	}
}"

jq '.S3PublishEndpoints = $config' --argjson config "$S3_CONFIG" < .aptly.conf > .aptly.new.conf
rm .aptly.conf
mv .aptly.new.conf .aptly.conf
