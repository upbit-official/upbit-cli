#!/usr/bin/env bash
# Auto-generated example — do not edit manually.
# Create Order (POST /v1/orders)
export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

# default
upbit orders create \
  --market "KRW-BTC" \
  --side "bid" \
  --volume "1" \
  --price "14000000" \
  --ord-type "limit"

# market_sell
upbit orders create \
  --market "KRW-BTC" \
  --side "ask" \
  --volume "0.001" \
  --ord-type "market"

# market_buy
upbit orders create \
  --market "KRW-BTC" \
  --side "bid" \
  --price "10000" \
  --ord-type "price"

# best_buy
upbit orders create \
  --market "KRW-BTC" \
  --side "bid" \
  --price "10000" \
  --ord-type "best" \
  --time-in-force "ioc"

# best_sell
upbit orders create \
  --market "KRW-BTC" \
  --side "ask" \
  --volume "0.001" \
  --ord-type "best" \
  --time-in-force "ioc"

# limit_ioc
upbit orders create \
  --market "KRW-BTC" \
  --side "bid" \
  --volume "0.001" \
  --price "10000000" \
  --ord-type "limit" \
  --time-in-force "ioc"

# limit_smp_reduce
upbit orders create \
  --market "KRW-BTC" \
  --side "bid" \
  --volume "0.001" \
  --price "10000000" \
  --ord-type "limit" \
  --smp-type "reduce"
