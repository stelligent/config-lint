#-------------------------------------------------------------------------------------------------------------
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License. See https://go.microsoft.com/fwlink/?linkid=2090316 for license information.
#-------------------------------------------------------------------------------------------------------------

FROM golang:1.13
ENV GO111MODULE=on

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

# Update package tool
RUN apt-get update

# Install apt-utils and suppress package configuration warning
RUN apt-get -y install --no-install-recommends apt-utils dialog 2>&1

# Install build tools
RUN apt-get install -y --no-install-recommends build-essential \
    git

# Install VS Code Go Dependencies
RUN go get -x -d github.com/stamblerre/gocode \
	&& go build -o gocode-gomod github.com/stamblerre/gocode \
	&& mv gocode-gomod $GOPATH/bin/ \
	&& go get -u -v \
	golang.org/x/tools/gopls \
	github.com/mdempsky/gocode \
	github.com/uudashr/gopkgs/cmd/gopkgs \
	github.com/sqs/goreturns \
	golang.org/x/lint/golint \
	github.com/ramya-rao-a/go-outline
