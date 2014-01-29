#!/bin/bash
set -e

# Run tests first
cd $(dirname $0)/../..
TEST_DIRS=(
	"command"
	"game"
	"game/acquire"
	"game/card"
	"game/farkle"
	"game/liars_dice"
	"game/log"
	"game/lost_cities"
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
