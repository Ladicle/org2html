#!/bin/sh

set -eu

IMAGE=golangci/golangci-lint:latest

USER_OVERRIDE=$(id -u ${USER}):$(id -g ${USER})
VOLUME=$(pwd):$(pwd)
WORKDIR=$(pwd)

docker run --rm -it -v "${VOLUME}" -w "${WORKDIR}" -u "${USER_OVERRIDE}" \
       -e GOLANGCI_LINT_CACHE=$(pwd)/.cache/golangci-lint \
       -e GOCACHE=$(pwd)/.cache/go-build \
       "${IMAGE}" \
       hack/golangci-lint-local.sh
