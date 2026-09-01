#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# List Transfers (GET /v1/pockets/transfers)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit pockets list-transfers \
  --direction "out" \
  --currency "XRP" \
  --limit 2
