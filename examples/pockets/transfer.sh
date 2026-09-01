#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# Transfer (POST /v1/pockets/transfers)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit pockets transfer \
  --to "9ca023a5-851b-4fec-9f0a-48cd83c2eaae" \
  --currency "XRP" \
  --amount "10.00" \
  --identifier "unique_transfer_id_20260531_01"
