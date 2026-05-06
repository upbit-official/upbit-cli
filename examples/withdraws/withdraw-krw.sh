#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# Withdraw KRW (POST /v1/withdraws/krw)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit withdraws create-krw-withdrawal \
  --amount "10000" \
  --two-factor-type "naver"
