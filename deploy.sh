#!/usr/bin/env bash

set -euo pipefail

PROJECT_ROOT="/home/kamiy2743/workspace/blog"
BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

NVM_BIN="/home/kamiy2743/.nvm/versions/node/v24.14.0/bin"

BACK_SERVICE_SRC="$BACKEND_DIR/service/blog-back.service"
BACK_SERVICE_DST="/etc/systemd/system/blog-back.service"
FRONT_SERVICE_SRC="$FRONTEND_DIR/service/blog-front.service"
FRONT_SERVICE_DST="/etc/systemd/system/blog-front.service"

cd "$FRONTEND_DIR"
PATH="$NVM_BIN:$PATH" npm install
PATH="$NVM_BIN:$PATH" npm run build

cd "$BACKEND_DIR"
/usr/local/go/bin/go build -o "$BACKEND_DIR/blog-back" ./cmd/blog

sudo rsync -av --delete "$BACK_SERVICE_SRC" "$BACK_SERVICE_DST"
sudo rsync -av --delete "$FRONT_SERVICE_SRC" "$FRONT_SERVICE_DST"
sudo chmod 644 "$BACK_SERVICE_DST" "$FRONT_SERVICE_DST"

sudo systemctl daemon-reload
sudo systemctl restart blog-back
sudo systemctl restart blog-front
