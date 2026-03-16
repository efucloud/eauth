#!/usr/bin/env bash

set -euo pipefail

REGISTRY="registry.cn-shenzhen.aliyuncs.com/efucloud-public/eauth"
GIT_COMMIT=$(git rev-parse HEAD)
BUILD_DATE=$(date +'%Y/%m/%d %H:%M:%S')
TAG="v1.0.0.$(date +'%Y%m%d%H%M')"

echo "Commit: $GIT_COMMIT"
echo "Build Date: $BUILD_DATE"
echo "Tag: $TAG"

# 构建并推送 amd64
docker buildx build \
  -f Dockerfile \
  --platform linux/amd64,linux/arm64 \
  --build-arg GIT_COMMIT="$GIT_COMMIT" \
  --build-arg BUILD_DATE="$BUILD_DATE" \
  -t "${REGISTRY}:${TAG}" \
  --push \
  .
#

echo "✅ Published multi-arch image: ${REGISTRY}:${TAG}"
