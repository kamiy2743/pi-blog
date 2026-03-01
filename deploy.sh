#!/usr/bin/env bash

set -euo pipefail

SERVICE_SRC="/home/kamiy2743/workspace/blog/service/blog.service"
SERVICE_DST="/etc/systemd/system/blog.service"

cd /home/kamiy2743/workspace/blog
/usr/local/go/bin/go build -o /home/kamiy2743/workspace/blog/blog ./cmd/blog

sudo rsync -av --delete $SERVICE_SRC $SERVICE_DST
sudo chmod 644 $SERVICE_DST

sudo systemctl daemon-reload
