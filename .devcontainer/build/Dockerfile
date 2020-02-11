FROM ubuntu:latest

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

# Install packages
RUN apt-get update \
		# Install apt-utils and suppress package configuration warning
		&& apt-get -y install --no-install-recommends apt-utils dialog 2>&1 \
		# Install build tools
		&& apt-get install -y \
		bc \
		unzip \
		wget \
		git \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
		ca-certificates \
		gnupg-agent \
		# Cleanup apt lists
		&& rm -rf /var/lib/apt/lists/*

# Install Go
ENV GOLANG_VERSION 1.13.7
RUN set -eux; \
	goRelArch='linux-amd64'; \
	goRelSha256='b3dd4bd781a0271b33168e627f7f43886b4c5d1c794a4015abf34e99c6526ca3'; \
	url="https://golang.org/dl/go${GOLANG_VERSION}.${goRelArch}.tar.gz"; \
	wget -O go.tgz "$url"; \
	echo "${goRelSha256} *go.tgz" | sha256sum -c -; \
	tar -C /usr/local -xzf go.tgz; \
	rm go.tgz
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
ENV GO111MODULE=on

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

# Create a non-root user
ARG USERNAME=config-lint-dev
ARG USER_UID=1000
ARG USER_GID=$USER_UID
RUN groupadd --gid $USER_GID $USERNAME \
  && useradd --uid $USER_UID --gid $USER_GID --shell /bin/bash -m $USERNAME

# Prompt gpg window inside container for signing commits & setup folder permissions for non-root user
RUN echo 'export GPG_TTY="$(tty)"' >> /home/$USERNAME/.bashrc \
  && mkdir /home/$USERNAME/.gnupg \
  && chown -R $USERNAME:$USERNAME /home/$USERNAME/.gnupg

# Persist bash history between runs
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
    && mkdir /commandhistory \
    && touch /commandhistory/.bash_history \
    && chown -R $USERNAME /commandhistory \
    && echo $SNIPPET >> "/home/$USERNAME/.bashrc"

# Add non-root user to $GOPATH
RUN chown -R $USERNAME $GOPATH \
		# Add write permission for /go/pkg
		&& chmod -R a+w /go/pkg

# Install Terraform
ARG TERRAFORM_VERSION=0.12.20
RUN cd /tmp \
	&& wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
	&& unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
	&& mv terraform /usr/local/bin/ \
	&& rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip

# Enter container as non-root user
USER $USERNAME
