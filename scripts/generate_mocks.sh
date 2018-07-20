#!/bin/bash

set -e

if [ -z "$GOPATH" ]; then
  echo "Please set GOPATH"
  exit 1
fi

mockgen_path=$GOPATH/bin/mockgen
base_package=github.com/activatedio/wrangle

if [ ! -f $mockgen_path ]; then
  echo "Installing gomock and mockgen"
  go get github.com/golang/mock/gomock
  go get github.com/golang/mock/mockgen
fi

generate() {
  package=$1
  mocks_dir=$package/mocks
  if [ ! -d $mocks_dir ]; then
    mkdir -p $mocks_dir
  fi
  find ./$package/*.go | xargs grep 'type \w* interface' | awk -vORS=, '{print $2}' | sed 's/,$//' | xargs -I % \
    $mockgen_path --package ${package}_mocks $base_package/$package %  > $mocks_dir/mocks.go
}

generate action
generate cloud_init 
generate image
generate network
generate server 
generate util
generate netboot
