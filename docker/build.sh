#!/bin/sh
# Build UpAndRunning2
# Script derived from https://github.com/gogits/gogs/blob/master/docker/build.sh (MIT-License)

set -x
set -e

# set path
export GOPATH=/tmp/go
export PATH=${PATH}:${GOPATH}/bin

# install dependencies
apk -U --no-progress add git go build-base ca-certificates

# build application
mkdir -p ${GOPATH}/src/github.com/MarvinMenzerath/
ln -s /app/upandrunning2/ ${GOPATH}/src/github.com/MarvinMenzerath/UpAndRunning2
cd ${GOPATH}/src/github.com/MarvinMenzerath/UpAndRunning2
go get
go build

# remove sources and unnecessary packages
rm -rf $GOPATH
apk --no-progress del git go build-base

# create custom user and group
addgroup -g 1777 uar2 && adduser -h /app/upandrunning2/ -H -D -G uar2 -u 1777 uar2

# set rights
chown -R uar2:uar2 /app/upandrunning2/