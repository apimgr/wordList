#!/usr/bin/env bash
set -euo pipefail

config_dir="${CONFIG_DIR:-/config/wordList}"
data_dir="${DATA_DIR:-/data/wordList}"
mkdir -p "$config_dir" "$data_dir"
chown wordlist:wordlist "$config_dir" "$data_dir"
exec su-exec wordlist /usr/local/bin/wordList \
  --config "$config_dir" \
  --data "$data_dir" \
  --address "${LISTEN_ADDRESS:-0.0.0.0}" \
  --port "${PORT:-80}" \
  --mode "${MODE:-production}" \
  "$@"
