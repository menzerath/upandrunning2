#!/bin/sh
# Build UpAndRunning2
# Script derived from https://github.com/gogits/gogs/blob/master/docker/build.sh (MIT-License)

set -x
set -e

export GOPATH=/tmp/go
export PATH=${PATH}:${GOPATH}/bin

apk -U --no-progress add git go build-base

mkdir -p ${GOPATH}/src/github.com/MarvinMenzerath/
ln -s /app/upandrunning2/ ${GOPATH}/src/github.com/MarvinMenzerath/UpAndRunning2
cd ${GOPATH}/src/github.com/MarvinMenzerath/UpAndRunning2
go get
go build

rm -rf $GOPATH
apk --no-progress del git go build-base

# keep this one to enable telegram notifications
apk -U --no-progress add ca-certificates