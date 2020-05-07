#!/bin/bash
if ! test -d dist
then
	echo "Missing dist directory from goreleaser step, exiting"
	exit 1
fi
TAG=$1
mv dist builds
tar czf $(basename ${GITHUB_REPOSITORY})_${TAG}.tgz .bonsai.yml builds/*.tar.gz builds/*_sha512-checksums.txt *.md
ls -l $(basename ${GITHUB_REPOSITORY})_${TAG}.tgz
