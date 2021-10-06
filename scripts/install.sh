#!/bin/bash

# install geth
wget -O /tmp/geth.tar.gz https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.10.9-eae3b194.tar.gz
tar -xvzf /tmp/geth.tar.gz
cp /tmp/geth-linux*/geth /usr/bin/geth
chmod +x /usr/bin/geth
