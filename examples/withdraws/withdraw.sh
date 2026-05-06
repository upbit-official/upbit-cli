#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# Withdraw Digital Asset (POST /v1/withdraws/coin)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit withdraws create-withdrawal \
  --currency "BTC" \
  --net-type "BTC" \
  --amount "0.01" \
  --address "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
