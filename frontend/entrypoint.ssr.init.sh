#!/usr/bin/env bash

set -euo pipefail

rm -rf ./dist/*
chown -R node:node ./dist/

exec su node /usr/local/bin/entrypoint.ssr.watch.sh
