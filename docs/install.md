# config-lint installation guide

## Homebrew

You can use [Homebrew](https://brew.sh/) to install the latest version:

``` bash
brew tap stelligent/tap
brew install config-lint
```

## Docker

You can pull the latest image from [DockerHub](https://hub.docker.com/r/stelligent/config-lint):

``` bash
docker pull stelligent/config-lint
```

If you choose to install and run via `docker` you will need mount a directory to the running container so that it has access to your configuration files.

``` bash
docker run -v /path/to/your/configs/:/foobar stelligent/config-lint -terraform /foobar/foo.tf

# or 

docker run --mount src=/path/to/your/configs/,target=/foobar,type=bind stelligent/config-lint -terraform /foobar/foo.tf
```

If you are linting Kubernetes configuration files, you will need to reference the path to the Kubernetes rules accordingly.

For example if the `pwd` has rules and configuration files:
```
docker run -v $(pwd):/foobar stelligent/config-lint -rules /foobar/path/to/rules/kubernetes.yml /foobar/path/to/configs
```

If you don't have your own set of rules that you want to run against your Kubernetes configuration files, you can copy or download the example set from [example-files/rules/kubernetes.yml](example-files/rules/kubernetes.yml).

## Linux

```
# Install the latest version of config-lint
curl -L https://github.com/stelligent/config-lint/releases/download/latest/config-lint_Linux_x86_64.tar.gz | tar xz -C /usr/local/bin config-lint

# See https://github.com/stelligent/config-lint/releases for release versions
VERSION=v1.0.0
curl -L https://github.com/stelligent/config-lint/releases/download/${VERSION}/config-lint_Linux_x86_64.tar.gz | tar xz -C /usr/local/bin config-lint

chmod +rx /usr/local/bin/config-lint
```

## Windows

Work in progress