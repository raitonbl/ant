#!/bin/bash
set -e

go test -short -coverprofile=bin/cov.out ./...

if sonar-scanner  > /dev/null
then
  sonar-scanner
fi