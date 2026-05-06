#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# Get Available Deposit Information (GET /v1/deposits/chance/coin)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit deposits retrieve-chance \
  --currency "BTC" \
  --net-type "BTC"
