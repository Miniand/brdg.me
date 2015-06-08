#!/bin/bash
set -e

DEPLOY_ADDRESS=${DEPLOY_ADDRESS:-root@api.brdg.me}

#Check we're 64-bit Linux
if [[ `uname` != "Linux" || `uname -m` != "x86_64" ]]; then
	echo "Can only deploy from 64-bit Linux"
	exit 1
fi

# Run tests first
cd $(dirname $0)
./test.sh

# Build and deploy files
cd ../main
godep go build
scp main $DEPLOY_ADDRESS:/usr/bin/brdg.me-new
ssh $DEPLOY_ADDRESS "if pgrep brdg.me; then service brdg.me stop; fi && mv -f /usr/bin/brdg.me-new /usr/bin/brdg.me && service brdg.me start"
rm main

echo "Deploy complete"
