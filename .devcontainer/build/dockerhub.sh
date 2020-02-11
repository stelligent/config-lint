#!/bin/bash -ex

set +x
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

# publish vscode-remote docker image to DockerHub, https://hub.docker.com/r/stelligent/vscode-remote-config-lint
docker build -t $DOCKER_ORG/vscode-remote-config-lint:${NEW_VERSION} --file .devcontainer/build/Dockerfile .
docker tag $DOCKER_ORG/vscode-remote-config-lint:${NEW_VERSION} $DOCKER_ORG/vscode-remote-config-lint:latest
docker push $DOCKER_ORG/vscode-remote-config-lint:${NEW_VERSION}
docker push $DOCKER_ORG/vscode-remote-config-lint:latest