#!/usr/bin/env bash

set -euo pipefail

rm -rf ./dist/*

npm run build:ssr:watch &

while [ ! -f ./dist/ssr/ssr.js ]; do
  sleep 1
done

exec node --watch ./dist/ssr/ssr.js
