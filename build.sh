#!/bin/bash

cd scanner
go mod tidy
docker buildx build --platform linux/amd64,linux/arm64 --tag back2basic/siaalert:latest -f ./docker/Dockerfile .
cd ..
