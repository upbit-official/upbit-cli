#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# Cancel Orders by IDs (DELETE /v1/orders/uuids)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit orders cancel-by-uuids \
  --uuid "bbbb8e07-1689-4769-af3e-a117016623f8" \
  --uuid "4312ba49-5f1a-4a01-9f3b-2d2bce17267e" \
  --uuid "bdb49a54-de36-4eb4-a963-9c8d4337a9da"
