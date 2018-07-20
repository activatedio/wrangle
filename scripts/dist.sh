#!/usr/bin/env bash
set -e
set -x

# Get the version from the environment, or try to figure it out from the build tags.
# We process the files in the same order Go does to find the last matching tag.
if [ -z $VERSION ]; then
    for file in $(ls version/version*.go | sort); do
        for tag in "$GOTAGS"; do
            if grep -q "// +build $tag" $file; then
                VERSION=$(awk -F\" '/Version =/ { print $2; exit }' <$file)
            fi
        done
    done
fi
if [ -z $VERSION ]; then
    echo "Please specify a version (couldn't find one based on build tags)."
    exit 1
fi
echo "==> Building version $VERSION..."

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change into that dir because we expect that.
cd $DIR

# Generate the tag.
if [ -z $NOTAG ]; then
  echo "==> Tagging..."
  git commit --allow-empty -a --gpg-sign=CF32B65A -m "Release v$VERSION"
  git tag -a -m "Version $VERSION" -s -u CF32B65A "v${VERSION}" master
fi

# Do a hermetic build inside a Docker container.
if [ -z $NOBUILD ]; then
    docker build -t activatedio/wrangle-builder scripts/wrangle-builder/
    docker run --rm -e "GOTAGS=$GOTAGS" -v "$(pwd)":/gopath/src/github.com/activatedio/wrangle activatedio/wrangle-builder ./scripts/dist_build.sh
fi

# Zip all the files.
rm -rf ./pkg/dist
mkdir -p ./pkg/dist
for FILENAME in $(find ./pkg -mindepth 1 -maxdepth 1 -type f); do
  FILENAME=$(basename $FILENAME)
  cp ./pkg/${FILENAME} ./pkg/dist/wrangle_${VERSION}_${FILENAME}
done

# Make the checksums.
pushd ./pkg/dist
shasum -a256 * > ./wrangle_${VERSION}_SHA256SUMS
if [ -z $NOSIGN ]; then
  echo "==> Signing..."
  gpg --default-key 348FFC4C --detach-sig ./wrangle_${VERSION}_SHA256SUMS
fi
popd

exit 0
