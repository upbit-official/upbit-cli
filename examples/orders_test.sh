#!/usr/bin/env bash
# Order Creation Test Scenario.
#
# Uses the Order Creation Test API to validate various order types without
# placing real orders.
#   - Check order availability
#   - Limit buy/sell test
#   - Market buy/sell test
#   - Best buy/sell test (IOC)
#   - Invalid order validation
#
# The Order Creation Test API performs the same validation as a real order,
# but no order is actually created -- it is fee-free and safe to use.
#
# Usage:
#   UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/orders_test.sh

set -euo pipefail

command -v jq >/dev/null 2>&1 || { echo "jq is required." >&2; exit 1; }
command -v bc >/dev/null 2>&1 || { echo "bc is required." >&2; exit 1; }

[ -n "${UPBIT_ACCESS_KEY+x}" ] && export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
[ -n "${UPBIT_SECRET_KEY+x}" ] && export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

MARKET="KRW-BTC"
TEST_VOLUME="0.0001"

divider() {
  echo ""
  printf '%0.s=' {1..60}; echo ""
  echo "  $1"
  printf '%0.s=' {1..60}; echo ""
}

# ---------------------------------------------------------------------------
# 1. Check order availability
# ---------------------------------------------------------------------------

section_check_order_chance() {
  divider "1. Order Availability"

  chance=$(upbit orders retrieve-chance --market "$MARKET")
  printf "  Market:          %s\n"   "$MARKET"
  printf "  Bid fee:         %s\n"   "$(echo "$chance" | jq -r '.bid_fee')"
  printf "  Ask fee:         %s\n"   "$(echo "$chance" | jq -r '.ask_fee')"
  printf "  Bid types:       %s\n"   "$(echo "$chance" | jq -r '.market.bid_types | @csv')"
  printf "  Ask types:       %s\n"   "$(echo "$chance" | jq -r '.market.ask_types | @csv')"
  printf "  Min bid amount:  %s\n"   "$(echo "$chance" | jq -r '.market.bid.min_total')"
  printf "  Min ask amount:  %s\n"   "$(echo "$chance" | jq -r '.market.ask.min_total')"
}

# ---------------------------------------------------------------------------
# 2. Limit buy/sell test
# ---------------------------------------------------------------------------

section_test_limit_orders() {
  divider "2. Limit Buy/Sell Test"

  tick_size=$(upbit orderbooks list-instruments --markets "$MARKET" | jq -r '.[0].tick_size')
  echo "  Tick size:      $tick_size"

  best_bid=$(upbit orderbooks list --markets "$MARKET" | jq -r '.[0].orderbook_units[0].bid_price')
  printf "  Best bid:       %.0f KRW\n" "$best_bid"

  bid_price=$(echo "scale=0; ($best_bid * 0.97) / $tick_size * $tick_size / 1" | bc)
  echo ""
  echo "  [Limit Buy]"
  printf "  Order price:    %.0f KRW (3%% discount)\n" "$bid_price"
  echo "  Order volume:   $TEST_VOLUME"
  bid_order=$(upbit orders test-create \
    --market "$MARKET" --side "bid" --ord-type "limit" \
    --price "$bid_price" --volume "$TEST_VOLUME")
  echo "  -> UUID:  $(echo "$bid_order" | jq -r '.uuid')"
  echo "  -> State: $(echo "$bid_order" | jq -r '.state')"
  echo "  -> Type:  $(echo "$bid_order" | jq -r '.ord_type')"

  ask_price=$(echo "scale=0; ($best_bid * 1.03) / $tick_size * $tick_size / 1" | bc)
  echo ""
  echo "  [Limit Sell]"
  printf "  Order price:    %.0f KRW (3%% premium)\n" "$ask_price"
  echo "  Order volume:   $TEST_VOLUME"
  ask_order=$(upbit orders test-create \
    --market "$MARKET" --side "ask" --ord-type "limit" \
    --price "$ask_price" --volume "$TEST_VOLUME")
  echo "  -> UUID:  $(echo "$ask_order" | jq -r '.uuid')"
  echo "  -> State: $(echo "$ask_order" | jq -r '.state')"
  echo "  -> Type:  $(echo "$ask_order" | jq -r '.ord_type')"
}

# ---------------------------------------------------------------------------
# 3. Market buy/sell test
# ---------------------------------------------------------------------------

section_test_market_orders() {
  divider "3. Market Buy/Sell Test"

  min_total=$(upbit orders retrieve-chance --market "$MARKET" | jq -r '.market.bid.min_total')

  echo "  [Market Buy]"
  printf "  Order amount:   %.0f KRW (minimum)\n" "$min_total"
  buy_order=$(upbit orders test-create \
    --market "$MARKET" --side "bid" --ord-type "price" --price "$min_total")
  echo "  -> UUID:  $(echo "$buy_order" | jq -r '.uuid')"
  echo "  -> State: $(echo "$buy_order" | jq -r '.state')"
  echo "  -> Type:  $(echo "$buy_order" | jq -r '.ord_type')"

  echo ""
  echo "  [Market Sell]"
  echo "  Order volume:   $TEST_VOLUME"
  sell_order=$(upbit orders test-create \
    --market "$MARKET" --side "ask" --ord-type "market" --volume "$TEST_VOLUME")
  echo "  -> UUID:  $(echo "$sell_order" | jq -r '.uuid')"
  echo "  -> State: $(echo "$sell_order" | jq -r '.state')"
  echo "  -> Type:  $(echo "$sell_order" | jq -r '.ord_type')"
}

# ---------------------------------------------------------------------------
# 4. Best buy/sell test (IOC)
# ---------------------------------------------------------------------------

section_test_best_orders() {
  divider "4. Best Buy/Sell Test (IOC)"

  min_total=$(upbit orders retrieve-chance --market "$MARKET" | jq -r '.market.bid.min_total')

  echo "  [Best Buy -- IOC]"
  printf "  Order amount:   %.0f KRW (minimum)\n" "$min_total"
  bid_order=$(upbit orders test-create \
    --market "$MARKET" --side "bid" --ord-type "best" \
    --price "$min_total" --time-in-force "ioc")
  echo "  -> UUID:  $(echo "$bid_order" | jq -r '.uuid')"
  echo "  -> State: $(echo "$bid_order" | jq -r '.state')"
  echo "  -> Type:  $(echo "$bid_order" | jq -r '.ord_type')"

  echo ""
  echo "  [Best Sell -- IOC]"
  echo "  Order volume:   $TEST_VOLUME"
  ask_order=$(upbit orders test-create \
    --market "$MARKET" --side "ask" --ord-type "best" \
    --volume "$TEST_VOLUME" --time-in-force "ioc")
  echo "  -> UUID:  $(echo "$ask_order" | jq -r '.uuid')"
  echo "  -> State: $(echo "$ask_order" | jq -r '.state')"
  echo "  -> Type:  $(echo "$ask_order" | jq -r '.ord_type')"
}

# ---------------------------------------------------------------------------
# 5. Invalid order validation
# ---------------------------------------------------------------------------

section_test_validation() {
  divider "5. Invalid Order Validation"

  echo "  [Validation] Non-existent market (KRW-INVALID)"
  if upbit orders test-create \
    --market "KRW-INVALID" --side "bid" --ord-type "limit" \
    --price "10000" --volume "1" 2>/dev/null; then
    echo "  -> Unexpectedly succeeded."
  else
    echo "  -> Error occurred (expected)"
  fi

  echo ""
  echo "  [Validation] Limit order without price"
  if upbit orders test-create \
    --market "$MARKET" --side "bid" --ord-type "limit" \
    --volume "$TEST_VOLUME" 2>/dev/null; then
    echo "  -> Unexpectedly succeeded."
  else
    echo "  -> Error occurred (expected)"
  fi
}

# ---------------------------------------------------------------------------
# Run
# ---------------------------------------------------------------------------

section_check_order_chance
section_test_limit_orders
section_test_market_orders
section_test_best_orders
section_test_validation

echo ""
printf '%0.s=' {1..60}; echo ""
echo "  Orders test scenario completed"
printf '%0.s=' {1..60}; echo ""
echo ""
echo "  * The Order Creation Test API does not create real orders."
echo "    Returned UUIDs cannot be used for lookups or cancellations."
