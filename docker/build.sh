#!/bin/sh
set -x
set -e

# Set temp environment vars
export GOPATH=/tmp/go
export PATH=${PATH}:${GOPATH}/bin

# Install build deps
apk -U --no-progress add git go

# Init go environment to build Gogs
mkdir -p ${GOPATH}/src/github.com/MarvinMenzerath/
ln -s /app/upandrunning2/ ${GOPATH}/src/github.com/MarvinMenzerath/UpAndRunning2
cd ${GOPATH}/src/github.com/MarvinMenzerath/UpAndRunning2
go get
go build

# Cleanup GOPATH
rm -rf $GOPATH

# Remove build deps
apk --no-progress del git go