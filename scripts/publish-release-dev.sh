#!/bin/bash

hash=$1
key=$2

curl -X POST -H "x-forta-key: $key" -H "content-type: application/json" -d "{\"hash\": \"$hash\"}" https://api.defender.openzeppelin.com/autotasks/47e42e49-4856-4e46-ad99-a6ef11f6324b/runs/webhook/62ea5767-415e-412d-aa34-ff31ed60b640/6QZzgKwrG2bX1qRc2ESprt