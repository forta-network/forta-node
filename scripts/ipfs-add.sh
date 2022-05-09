#!/bin/bash


curl -X POST -F file=@$1 https://ipfs.forta.network/api/v0/add
