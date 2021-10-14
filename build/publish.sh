#!/bin/bash

ver=`cat docs/version`
aws s3 cp ant s3://${{ secrets.AWS_BUCKET }}/repositories/binaries/ant/$VERSION/$S3_DIRECTORY/application
