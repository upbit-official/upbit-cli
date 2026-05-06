#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# List Orders by IDs (GET /v1/orders/uuids)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit orders list-by-uuids \
  --uuid "5d303952-8be9-41e6-915b-121a90026248" \
  --uuid "3944c2c1-bd8c-441a-aa25-2370d08217a9" \
  --uuid "5b95451b-971e-4e76-8f61-5ff441f078d5" \
  --uuid "3b67e543-8ad3-48d0-8451-0dad315cae73"
