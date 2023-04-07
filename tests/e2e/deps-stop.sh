#!/bin/bash

sudo pkill geth
sudo pkill ipfs
sudo pkill disco
kill "$(pidof mock-graphql-api)"

