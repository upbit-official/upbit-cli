#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# Create Deposit Address (POST /v1/deposits/generate_coin_address)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit deposits create-coin-address \
  --currency "BTC" \
  --net-type "BTC"
