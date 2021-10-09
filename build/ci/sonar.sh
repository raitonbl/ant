#!/bin/bash

go test -short -coverprofile=bin/cov.out ./...
sonar-scanner