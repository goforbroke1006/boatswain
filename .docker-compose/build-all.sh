#!/bin/bash

set -e

SCRIPT_PATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

(
  cd "${SCRIPT_PATH}"/..

  IMAGE_NAME=$(cat go.mod | grep '^module ' | awk '{print $2}')
  PROJECT_NAME="${IMAGE_NAME##*/}"

  docker build \
    -f ./.docker-compose/"${PROJECT_NAME}"/Dockerfile \
    -t "${IMAGE_NAME}":dev \
    --build-arg PROJECT_NAME="${PROJECT_NAME}" \
    ./
)
