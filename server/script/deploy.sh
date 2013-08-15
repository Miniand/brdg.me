#!/bin/bash
set -e

DEPLOY_ADDRESS=${DEPLOY_ADDRESS:-root@brdg.me}

#Check we're 64-bit Linux
if [[ `uname` != "Linux" || `uname -m` != "x86_64" ]]; then
	echo "Can only deploy from 64-bit Linux"
	exit 1
fi

# Run tests first
cd $(dirname $0)
./test.sh

# Build and deploy files
cd ../email
go get
go build
ssh $DEPLOY_ADDRESS service brdg.me-email stop
scp email $DEPLOY_ADDRESS:/usr/bin/brdg.me-email
ssh $DEPLOY_ADDRESS service brdg.me-email start
rm email
cd ../web
go get
go build
ssh $DEPLOY_ADDRESS service brdg.me-web stop
scp web $DEPLOY_ADDRESS:/usr/bin/brdg.me-web
ssh $DEPLOY_ADDRESS service brdg.me-web start
rm web

echo "Deploy complete"
