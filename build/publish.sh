#!/bin/bash

ver=`cat docs/version`
aws s3 cp ant s3://$SECRET_AWS_BUCKET/repositories/binaries/ant/$ver/$S3_DIRECTORY/application
