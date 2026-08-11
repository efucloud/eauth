#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
WEB_DIR="${ROOT_DIR}/web"
WEB_DIST_DIR="${WEB_DIR}/dist"
EMBED_WEB_DIR="${ROOT_DIR}/pkg/embeds/web"
EMBED_WEB_TMP_DIR="${ROOT_DIR}/pkg/embeds/.web.tmp"

if command -v corepack >/dev/null 2>&1; then
  corepack enable >/dev/null 2>&1 || true
fi

cd "${WEB_DIR}"
HUSKY=0 yarn install --frozen-lockfile --non-interactive
yarn build

if [[ ! -f "${WEB_DIST_DIR}/index.html" ]]; then
  echo "frontend build did not produce ${WEB_DIST_DIR}/index.html" >&2
  exit 1
fi

rm -rf "${EMBED_WEB_TMP_DIR}"
mkdir -p "${EMBED_WEB_TMP_DIR}"
cp -R "${WEB_DIST_DIR}/." "${EMBED_WEB_TMP_DIR}/"
rm -rf "${EMBED_WEB_DIR}"
mv "${EMBED_WEB_TMP_DIR}" "${EMBED_WEB_DIR}"

echo "Embedded frontend assets into ${EMBED_WEB_DIR}"
