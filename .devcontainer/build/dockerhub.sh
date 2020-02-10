#!/bin/bash -ex

set +x
if [[ -z ${DOCKER_USER} ]];
then
  echo DOCKER_USER must be set in the environment
  exit 1
fi
if [[ -z ${DOCKER_PASSWORD} ]];
then
  echo DOCKER_PASSWORD must be set in the environment
  exit 1
fi
if [[ -z ${DOCKER_ORG} ]];
then
  echo DOCKER_ORG must be set in the environment
  exit 1
fi
if [[ -z ${NEW_VERSION} ]];
then
  echo NEW_VERSION must be set in the environment
  exit 1
fi
set -x

# Docker Login
set +x
echo $DOCKER_PASSWORD | docker login -u $DOCKER_USER --password-stdin
set -x

# publish vscode-remote docker image to DockerHub, https://hub.docker.com/r/stelligent/vscode-remote-config-lint
docker build -t $DOCKER_ORG/vscode-remote-config-lint:${NEW_VERSION} --file .devcontainer/build/Dockerfile .
docker tag $DOCKER_ORG/vscode-remote-config-lint:${NEW_VERSION} $DOCKER_ORG/vscode-remote-config-lint:latest
docker push $DOCKER_ORG/vscode-remote-config-lint:${NEW_VERSION}
docker push $DOCKER_ORG/vscode-remote-config-lint:latest