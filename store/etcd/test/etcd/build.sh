#!/usr/bin/env bash

export GOARCH=amd64
export GOOS=linux

go build -o etcdTestCli
scp etcdTestCli nw@cnw4:/home/nw/testing/grpcld
