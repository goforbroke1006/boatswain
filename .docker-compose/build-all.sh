#!/bin/bash

set -e

(
  cd ..

  PROJECT_NAME=$(basename "$(pwd)")

  docker build \
    -f ./.docker-compose/"${PROJECT_NAME}"/Dockerfile \
    -t local.env/"${PROJECT_NAME}":dev \
    ./
)
