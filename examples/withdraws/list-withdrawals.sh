#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# List Withdrawals (GET /v1/withdraws)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit withdraws list \
  --currency "XRP" \
  --state "DONE"
