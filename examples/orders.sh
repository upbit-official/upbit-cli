#!/usr/bin/env bash
# Order Creation and Management Scenario.
#
# Demonstrates how to create and manage orders using the Upbit Exchange API.
#   - Check order availability
#   - Place a limit buy order -> look up -> cancel
#   - Check holdings + verify market sell eligibility
#   - List completed orders
#
# Dry run (default):
#   Read-only lookups. No orders are placed.
#
# Live run (DRY_RUN=false):
#   Places a limit buy, looks it up, then cancels it.
#
# Usage:
#   UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/orders.sh
#   DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/orders.sh

set -euo pipefail

command -v jq >/dev/null 2>&1 || { echo "jq is required." >&2; exit 1; }
command -v bc >/dev/null 2>&1 || { echo "bc is required." >&2; exit 1; }

[ -n "${UPBIT_ACCESS_KEY+x}" ] && export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
[ -n "${UPBIT_SECRET_KEY+x}" ] && export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

MARKET="KRW-BTC"
DRY_RUN="${DRY_RUN:-true}"

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
# 2. Limit buy -> look up -> cancel
# ---------------------------------------------------------------------------

section_limit_bid_and_cancel() {
  divider "2. Limit Buy -> Look Up -> Cancel"

  tick_size=$(upbit orderbooks list-instruments --markets "$MARKET" | jq -r '.[0].tick_size')
  echo "  Tick size:      $tick_size"

  best_bid=$(upbit orderbooks list --markets "$MARKET" | jq -r '.[0].orderbook_units[0].bid_price')
  printf "  Best bid:       %.0f KRW\n" "$best_bid"

  target_price=$(echo "scale=0; ($best_bid * 0.97) / $tick_size * $tick_size / 1" | bc)
  volume="0.0001"
  printf "  Order price:    %'.0f KRW (3%% discount)\n" "$target_price"
  echo "  Order volume:   $volume"

  if [ "$DRY_RUN" != "false" ]; then
    echo ""
    echo "  [Skipped] Dry run mode -- actual order omitted."
    return 0
  fi

  order=$(upbit orders create \
    --market "$MARKET" \
    --side "bid" \
    --ord-type "limit" \
    --price "$target_price" \
    --volume "$volume")
  uuid=$(echo "$order" | jq -r '.uuid')
  echo ""
  echo "  Order created"
  echo "  UUID:   $uuid"
  echo "  State:  $(echo "$order" | jq -r '.state')"

  info=$(upbit orders retrieve --uuid "$uuid")
  echo "  Lookup: state=$(echo "$info" | jq -r '.state'), executed=$(echo "$info" | jq -r '.executed_volume')"

  upbit orders cancel --uuid "$uuid" > /dev/null
  for _ in $(seq 1 10); do
    cancelled=$(upbit orders retrieve --uuid "$uuid")
    state=$(echo "$cancelled" | jq -r '.state')
    if [ "$state" = "done" ] || [ "$state" = "cancel" ]; then
      break
    fi
    sleep 0.5
  done
  echo "  Cancel: state=$state"
}

# ---------------------------------------------------------------------------
# 3. Check holdings + market sell eligibility
# ---------------------------------------------------------------------------

section_check_market_sell() {
  divider "3. Holdings + Market Sell Eligibility"

  coin=$(echo "$MARKET" | cut -d'-' -f2)

  accounts=$(upbit accounts list)
  target=$(echo "$accounts" | jq --arg c "$coin" '.[] | select(.currency == $c)')

  if [ -z "$target" ] || [ "$(echo "$target" | jq -r '.balance')" = "0" ]; then
    echo "  No $coin holdings."
  else
    balance=$(echo "$target" | jq -r '.balance')
    avg_buy=$(echo "$target" | jq -r '.avg_buy_price')
    echo "  Asset:          $coin"
    printf "  Balance:        %.8f\n" "$balance"
    printf "  Avg buy price:  %.0f KRW\n" "$avg_buy"
  fi

  chance=$(upbit orders retrieve-chance --market "$MARKET")
  ask_types=$(echo "$chance" | jq -r '.market.ask_types | @csv')
  if echo "$ask_types" | grep -q "market"; then
    echo "  Market sell:    supported"
  else
    echo "  Market sell:    not supported (ask_types=$ask_types)"
  fi
}

# ---------------------------------------------------------------------------
# 4. List completed orders
# ---------------------------------------------------------------------------

section_list_closed_orders() {
  divider "4. Completed Orders"

  closed=$(upbit orders list-closed --market "$MARKET" --limit 5)
  count=$(echo "$closed" | jq 'length')

  if [ "$count" = "0" ]; then
    echo "  No completed orders."
    return 0
  fi

  printf "  %-10s %4s %8s %6s %14s\n" "UUID" "Side" "Type" "State" "Executed"
  echo "  --------------------------------------------------"
  echo "$closed" | jq -r '.[] | [.uuid[:8], .side, .ord_type, .state, .executed_volume] | @tsv' | \
    awk -F'\t' '{
      side = ($2 == "bid") ? "BUY" : "SELL"
      printf "  %s... %4s %8s %6s %14s\n", $1, side, $3, $4, $5
    }'
}

# ---------------------------------------------------------------------------
# Run
# ---------------------------------------------------------------------------

section_check_order_chance || true
section_limit_bid_and_cancel || true
section_check_market_sell || true
section_list_closed_orders || true

echo ""
printf '%0.s=' {1..60}; echo ""
echo "  Orders scenario completed"
printf '%0.s=' {1..60}; echo ""
echo ""
echo "  * This example is for educational purposes only."
echo "    Use at your own risk for actual trading."
