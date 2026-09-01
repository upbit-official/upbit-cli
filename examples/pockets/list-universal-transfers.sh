#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# List Universal Transfers (GET /v1/pockets/universal_transfers)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit pockets list-universal-transfers \
  --state "processing" \
  --state "done" \
  --currency "XRP" \
  --limit 2
