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

# Create a non-root user and provide sudo rights
ARG USERNAME=config-lint-dev
ARG USER_UID=1000
ARG USER_GID=$USER_UID
RUN groupadd --gid $USER_GID $USERNAME \
  && useradd --uid $USER_UID --gid $USER_GID --shell /bin/bash -m $USERNAME

# Prompt gpg window inside container for signing commits & setup folder permissions for non-root user
RUN echo 'export GPG_TTY="$(tty)"' >> /home/$USERNAME/.bashrc \
  && mkdir /home/$USERNAME/.gnupg \
  && chown -R $USERNAME /home/$USERNAME/.gnupg

# Enter container as non-root user
USER $USERNAME
