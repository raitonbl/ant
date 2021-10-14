#!/bin/bash
set -e

ver=`cat docs/version`
ver="$ver-$GITHUB_LABEL.$GITHUB_RUN_ID.$GITHUB_RUN_NUMBER"

echo "$ver" > "docs/version"

go build
