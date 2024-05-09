#!/bin/bash

# Change the current directory to the directory of the script
cd "$(dirname "$0")"

go mod init list && go get
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go
zip myfunc.zip bootstrap
rm -rf bootstrap

aws.exe lambda update-function-code --function-name listTweet --zip-file fileb://myfunc.zip