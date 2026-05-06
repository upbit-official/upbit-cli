#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# Get Order (GET /v1/order)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit orders retrieve \
  --uuid "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"

# by_identifier
upbit orders retrieve \
  --identifier "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
