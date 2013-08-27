#!/bin/bash
set -e

# Run tests first
cd $(dirname $0)/../..
TEST_DIRS=(
	"command"
	"game"
	"game/card"
	"game/farkle"
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