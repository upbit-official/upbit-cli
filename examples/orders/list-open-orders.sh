#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# List Open Orders (GET /v1/orders/open)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit orders list-open \
  --market "KRW-BTC" \
  --state "wait" \
  --state "watch"
