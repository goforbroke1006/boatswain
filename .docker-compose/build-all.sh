#!/bin/bash

set -e

SCRIPT_PATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

(
  cd "${SCRIPT_PATH}"/..

  PROJECT_NAME=$(basename "$(pwd)")

  docker build \
    -f ./.docker-compose/"${PROJECT_NAME}"/Dockerfile \
    -t local.env/"${PROJECT_NAME}":dev \
    ./
)
