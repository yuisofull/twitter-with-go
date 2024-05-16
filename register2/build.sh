#!/bin/bash

# Change the current directory to the directory of the script
cd "$(dirname "$0")"

GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go
zip myfunc.zip bootstrap
rm -rf bootstrap

aws.exe lambda update-function-code --function-name register2 --zip-file fileb://myfunc.zip