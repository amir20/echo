#!/bin/bash
set -e

IMAGE="amir20/echo"
TAG="${1:-latest}"

docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --tag "${IMAGE}:${TAG}" \
  --push \
  .
