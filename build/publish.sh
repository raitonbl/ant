#!/bin/bash

ver=`cat docs/version`

if [ -z "${S3_FILE}" ]; then
    filename="ant"
else
    filename=$S3_FILE
fi

aws s3 cp $filename s3://$SECRET_AWS_BUCKET/repositories/binaries/ant/$ver/$S3_DIRECTORY/$filename
