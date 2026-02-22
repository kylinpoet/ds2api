#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT_DIR"

./tests/scripts/check-node-split-syntax.sh
node --test api/helpers/stream-tool-sieve.test.js api/chat-stream.test.js api/compat/js_compat_test.js "$@"
