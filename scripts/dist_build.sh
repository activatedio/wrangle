#!/usr/bin/env bash
set -e

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that dir because we expect that.
cd $DIR

# Make sure build tools are available.
make tools

# Fixup the timestamps to match what's checked in. This will allow us to cleanly
# verify that the checked-in content is up to date without spurious diffs of the
# file mod times.
pushd pkg
cat ../command/agent/bindata_assetfs.go | ../scripts/fixup_times.sh
popd

# Now we are ready to do a clean build of everything. The "all" build will blow
# away our pkg folder so we have to regenerate the ui once more. This is probably
# for the best since we have meddled with the timestamps.
rm -rf pkg
make all
chmod 777 pkg

exit 0
