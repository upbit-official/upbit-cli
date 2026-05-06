#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# Verify Travel Rule by Deposit TxID (POST /v1/travel_rule/deposit/txid)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit travel-rule verify-deposit-by-txid \
  --vasp-uuid "8d4fe968-82b2-42e5-822f-3840a245f802" \
  --txid "5b871d34-fe38-4025-8f5c-9b22028f85d3" \
  --currency "ETH" \
  --net-type "ETH"
