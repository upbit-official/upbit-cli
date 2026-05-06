#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# Deposit KRW (POST /v1/deposits/krw)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit deposits deposit-krw \
  --amount "10000" \
  --two-factor-type "naver"
