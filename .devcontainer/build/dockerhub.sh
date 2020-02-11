#!/bin/bash -ex

set +x
if [[ -z ${DOCKER_ORG} ]];
then
  echo DOCKER_ORG must be set in the environment
  exit 1
fi
if [[ -z ${GITHUB_SHA} ]];
then
  echo GITHUB_SHA must be set in the environment
  exit 1
fi
set -x

COMMIT_HASH=${GITHUB_SHA:0:8}

# publish vscode-remote docker image to DockerHub, https://hub.docker.com/r/stelligent/vscode-remote-config-lint
docker build -t $DOCKER_ORG/vscode-remote-config-lint:${COMMIT_HASH} --file .devcontainer/build/Dockerfile .
docker tag $DOCKER_ORG/vscode-remote-config-lint:${COMMIT_HASH} $DOCKER_ORG/vscode-remote-config-lint:latest
docker push $DOCKER_ORG/vscode-remote-config-lint:${COMMIT_HASH}
docker push $DOCKER_ORG/vscode-remote-config-lint:latest