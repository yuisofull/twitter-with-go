#!/bin/bash

go mod init create && go get
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap main.go
zip myfunc.zip bootstrap
rm -rf bootstrap