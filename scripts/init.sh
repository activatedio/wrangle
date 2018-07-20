#!/usr/bin/env bash
#
# This script builds the application from source for multiple platforms.
set -e

export CGO_ENABLED=0

go get github.com/mitchellh/gox

