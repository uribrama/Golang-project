#!/bin/bash

aws lambda create-function \
        --function-name persons_function \
        --runtime go1.x \
        --handler persons \
        --role arn:aws:iam::283132657931:role/lambda-ex  \
        --zip-file fileb://persons.zip \
        --region us-east-2 \
        --timeout 10      \
        --environment Variables="{DB_HOST=${DB_HOST},           \
                              DB_NAME=${DB_NAME},   \
                              DB_PASSWORD=${DB_PASSWORD},   \
                              DB_PORT=${DB_PORT}}",  \
                              DB_USER=${DB_USER}}",  \
                              GO_ENV=${testing}}" 
