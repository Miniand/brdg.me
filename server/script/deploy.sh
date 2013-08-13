#!/bin/sh
set -e

DEPLOY_ADDRESS=${DEPLOY_ADDRESS:-root@brdg.me}

#Check we're 64-bit Linux
platform=`uname`
architecture=`uname -m`
if [[ `uname` != "Linux" || `uname -m` != "x86_64" ]]; then
	echo "Can only deploy from 64-bit Linux"
	exit 1
fi

# Run tests first
cd $(dirname $0)/../..
TEST_DIRS=(
	"command"
	"game"
	"game/card"
	"game/log"
	"game/poker"
	"game/no_thanks"
	"game/tic_tac_toe"
	"render"
	"server/email"
	"server/model"
)
for i in "${TEST_DIRS[@]}"
do
	go test -i ./$i
	go test ./$i
done

# Build and deploy files
cd server/email
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