#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# Batch Cancel Orders (DELETE /v1/orders/open)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit orders cancel-open \
  --quote-currencies "KRW,BTC" \
  --cancel-side "all" \
  --excluded-pairs "KRW-ETH,BTC-XRP"
