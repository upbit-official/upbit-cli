#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# List Closed Orders (GET /v1/orders/closed)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit orders list-closed \
  --market "KRW-BTC" \
  --state "done" \
  --state "cancel"
