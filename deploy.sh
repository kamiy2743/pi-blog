#!/usr/bin/env bash

set -euo pipefail

cd /home/kamiy2743/workspace/blog
/usr/local/go/bin/go build -o /home/kamiy2743/workspace/blog/blog ./cmd/blog
