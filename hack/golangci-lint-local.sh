#!/bin/sh

set -ex

golangci-lint run -c .golangci.yml org/...
