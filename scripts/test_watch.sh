#!/bin/bash

pattern=$1
if [ -z "$pattern" ]; then
  pattern="./..."
fi

rerun --pattern="**/*.go" "go test $pattern"
