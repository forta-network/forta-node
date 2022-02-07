#!/bin/sh

set -x

REVISION="$1"
SEMVER="$2"
GPG_NAME="$3"
GPG_PASSPHRASE="$4"

# remove old package so we don't get an error if we try to add the same
DEBIAN_VERSION=$(./scripts/debian-version.sh "$SEMVER")
aptly repo show -with-packages forta
aptly repo remove forta 'forta_'"$DEBIAN_VERSION"'_amd64'

set -e

aptly repo add forta apt/forta_*_amd64.deb
aptly snapshot create "$REVISION" from repo forta
# "-batch=true" flag here makes passphrase work: https://github.com/aptly-dev/aptly/issues/642
aptly publish snapshot -distribution=stable -batch=true -force-overwrite -gpg-key="$GPG_NAME" \
	-passphrase="$GPG_PASSPHRASE" "$REVISION" s3:releaseBucket:repositories/apt/
