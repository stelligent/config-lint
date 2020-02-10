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
		bc \
		unzip \
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
  && useradd --uid $USER_UID --gid $USER_GID --shell /bin/bash -m $USERNAME \
	&& apt-get install -y sudo \
	&& echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
	&& chmod 0440 /etc/sudoers.d/$USERNAME

# Prompt gpg window inside container for signing commits & setup folder permissions for non-root user
RUN echo 'export GPG_TTY="$(tty)"' >> /home/$USERNAME/.bashrc \
  && mkdir /home/$USERNAME/.gnupg \
  && chown -R $USERNAME /home/$USERNAME/.gnupg

# Persist bash history between runs
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
    && mkdir /commandhistory \
    && touch /commandhistory/.bash_history \
    && chown -R $USERNAME /commandhistory \
    && echo $SNIPPET >> "/home/$USERNAME/.bashrc"

# Add non-root user to $GOPATH
RUN chown -R $USERNAME $GOPATH

# Add write permission for /go/pkg
RUN chmod -R a+w /go/pkg

# Install Terraform
ARG TERRAFORM_VERSION=0.12.20
RUN cd /tmp \
	&& wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
	&& unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
	&& mv terraform /usr/local/bin/

# Enter container as non-root user
USER $USERNAME
