#!/bin/bash

# Build Script for creating multiple architecture releases

# Requires:
# go get github.com/mitchellh/gox

## get version from self to include in file names
go build
VERSION=`hastie --version | sed -e 's/hastie v//'`

echo "Building $VERSION"
gox -osarch="linux/amd64 darwin/amd64 windows/amd64" -output "{{.Dir}}-$VERSION-{{.OS}}/{{.Dir}}"

for arch in linux darwin windows; do
    tar cf hastie-$VERSION-$arch.tar hastie-$VERSION-$arch
    gzip hastie-$VERSION-$arch.tar
    rm -rf hastie-$VERSION-$arch
done

