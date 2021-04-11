#! /usr/bin/env sh

set -e

# build the docker file
GIT_COMMIT=$(git rev-list -1 HEAD) && \
DOCKER_BUILDKIT=1 docker pull keratin/authn-server:latest
