#!/bin/bash

apt install curl

# Build the project
make

# Build GO App
export PATH=$PATH:/usr/local/go/bin
cd ./user
go get .
go build .

# Check if build was successful
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi
