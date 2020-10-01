#!/bin/bash

## requires go 1.13+

export PATH=/usr/local/go/bin:$PATH

if [ ! -z "$debugBuild" ]; then export BUILDFLAGS='-gcflags "all=-N -l"'; fi
go build -ldflags '-linkmode=external'  -v -o pac-sliprovider ./