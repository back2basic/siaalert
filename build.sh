#!/bin/bash

cd siaalert
go mod tidy
docker buildx build --platform linux/amd64,linux/arm64 --tag back2basic/siaalert:latest -f ./docker/Dockerfile --push .
cd ..
