#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# Test Order Creation (POST /v1/orders/test)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit orders test-create \
  --market "KRW-BTC" \
  --side "bid" \
  --volume "1" \
  --price "14000000" \
  --ord-type "limit"
