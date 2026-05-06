#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# Cancel and New Order (POST /v1/orders/cancel_and_new)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit orders cancel-and-new \
  --prev-order-uuid "ad217e24-ed02-469c-9b30-c08dbbda6908" \
  --new-ord-type "limit" \
  --new-price "100000000"
