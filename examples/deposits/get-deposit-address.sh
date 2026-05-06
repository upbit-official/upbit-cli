#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# Get Deposit Address (GET /v1/deposits/coin_address)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit deposits retrieve-coin-address \
  --currency "BTC" \
  --net-type "BTC"
