#!/bin/sh
set -eu

VERSION=${1:?usage: ./update.sh <pi-version>}
BASE_URL="https://pi.dev/api/installer/releases/${VERSION}"
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

curl -fsSL "${BASE_URL}/package.json" -o "${TMP_DIR}/package.json"
curl -fsSL "${BASE_URL}/package-lock.json" -o "${TMP_DIR}/package-lock.json"

mv "${TMP_DIR}/package.json" ./package.json
mv "${TMP_DIR}/package-lock.json" ./package-lock.json

printf 'Updated Pi manifests to %s\n' "$VERSION"
