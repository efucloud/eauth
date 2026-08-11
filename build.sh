#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
IMAGE_REPO="${IMAGE_REPO:-registry.cn-shenzhen.aliyuncs.com/efucloud-public/eauth}"
IMAGE_TAG="${IMAGE_TAG:-v1.0.0.$(date +'%Y%m%d%H%M')}"
PLATFORMS="${PLATFORMS:-linux/amd64,linux/arm64}"
DOCKERFILE="${DOCKERFILE:-Dockerfile.local}"
DOCKER_BUILD_FLAGS="${DOCKER_BUILD_FLAGS:-}"
GIT_COMMIT="$(git -C "$ROOT_DIR" rev-parse HEAD 2>/dev/null || printf unknown)"
BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
IMAGE="${IMAGE_REPO}:${IMAGE_TAG}"

if ! command -v docker >/dev/null 2>&1; then
  echo "docker command not found"
  exit 1
fi

echo "Root dir: ${ROOT_DIR}"
echo "Image repo: ${IMAGE_REPO}"
echo "Image tag: ${IMAGE_TAG}"
echo "Platforms: ${PLATFORMS}"
echo "Dockerfile: ${DOCKERFILE}"
echo "Git commit: ${GIT_COMMIT}"
echo "Build date: ${BUILD_DATE}"

docker buildx build \
  -f "${ROOT_DIR}/${DOCKERFILE}" \
  --platform "${PLATFORMS}" \
  --build-arg GIT_COMMIT="${GIT_COMMIT}" \
  --build-arg BUILD_DATE="${BUILD_DATE}" \
  ${DOCKER_BUILD_FLAGS} \
  -t "${IMAGE}" \
  --provenance=false \
  --progress=plain \
  --push \
  "${ROOT_DIR}"

echo "Build and push complete: ${IMAGE}"
