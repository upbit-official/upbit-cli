#!/usr/bin/env bash
# DCA (Dollar Cost Averaging) Automated Market Buy Scenario.
#
# Demonstrates a recurring fixed-amount market buy using the Upbit Exchange API.
#   - Fetch current price and verify buy configuration
#   - Repeat market buy for the specified number of rounds
#   - Summarize buy results (total spent, total volume, average price)
#
# Dry run (default):
#   Fetches the current price only. No orders are placed.
#
# Live run (DRY_RUN=false):
#   Performs real market-buy orders after verifying the balance.
#
# Usage:
#   UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/dca.sh
#   DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/dca.sh

set -euo pipefail

command -v jq >/dev/null 2>&1 || { echo "jq is required." >&2; exit 1; }
command -v bc >/dev/null 2>&1 || { echo "bc is required." >&2; exit 1; }

[ -n "${UPBIT_ACCESS_KEY+x}" ] && export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
[ -n "${UPBIT_SECRET_KEY+x}" ] && export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

MARKET="KRW-BTC"
BUY_AMOUNT="5000"   # Amount per buy (KRW)
TOTAL_ROUNDS=3      # Total number of buys
INTERVAL=5          # Interval between buys (seconds)
DRY_RUN="${DRY_RUN:-true}"

divider() {
  echo ""
  printf '%0.s=' {1..60}; echo ""
  echo "  $1"
  printf '%0.s=' {1..60}; echo ""
}

current_price() {
  upbit tickers list-by-trading-pairs --markets "$MARKET" | jq -r '.[0].trade_price'
}

# ---------------------------------------------------------------------------
# 1. DCA configuration
# ---------------------------------------------------------------------------

section_setup() {
  local mode_label
  [ "$DRY_RUN" = "false" ] && mode_label="LIVE RUN" || mode_label="DRY RUN"
  divider "1. DCA Configuration [$mode_label]"

  local current
  current=$(current_price)
  local total_needed
  total_needed=$(echo "$BUY_AMOUNT * $TOTAL_ROUNDS" | bc)

  if [ "$DRY_RUN" = "false" ]; then
    fiat=$(echo "$MARKET" | cut -d'-' -f1)
    accounts=$(upbit accounts list)
    available=$(echo "$accounts" | jq --arg c "$fiat" '.[] | select(.currency == $c) | .balance' | tr -d '"')
    available=${available:-0}

    if (( $(echo "$available < $total_needed" | bc -l) )); then
      echo "  [Error] $fiat balance ($available) is less than required ($total_needed)." >&2
      exit 1
    fi

    printf "  Current price:    %.0f KRW\n" "$current"
    printf "  %s balance:       %.0f KRW\n" "$fiat" "$available"
  else
    printf "  Current price:    %.0f KRW\n" "$current"
  fi

  echo "  Market:           $MARKET"
  printf "  Amount per buy:   %.0f KRW\n" "$BUY_AMOUNT"
  echo "  Rounds:           $TOTAL_ROUNDS"
  echo "  Interval:         ${INTERVAL}s"
  printf "  Total required:   %.0f KRW\n" "$total_needed"
  [ "$DRY_RUN" != "false" ] && echo "" && echo "  * Dry run mode: no orders are placed."
}

# ---------------------------------------------------------------------------
# 2. Execute DCA
# ---------------------------------------------------------------------------

section_execute_dca() {
  divider "2. Execute DCA"

  total_spent=0
  total_volume=0
  success_count=0

  if [ "$DRY_RUN" != "false" ]; then
    echo "  [Skipped] Dry run mode -- actual orders omitted."
    echo ""
    for r in $(seq 1 "$TOTAL_ROUNDS"); do
      current=$(current_price)
      est_volume=$(echo "scale=8; $BUY_AMOUNT / $current" | bc)
      printf "  [%d/%d] Price: %.0f KRW  (expected volume: %s)\n" "$r" "$TOTAL_ROUNDS" "$current" "$est_volume"
      [ "$r" -lt "$TOTAL_ROUNDS" ] && sleep "$INTERVAL"
    done
    return 0
  fi

  for r in $(seq 1 "$TOTAL_ROUNDS"); do
    printf "\n  [%d/%d] Market buy (%s KRW) ...\n" "$r" "$TOTAL_ROUNDS" "$BUY_AMOUNT"

    order=$(upbit orders create \
      --market "$MARKET" --side "bid" --ord-type "price" --price "$BUY_AMOUNT")
    uuid=$(echo "$order" | jq -r '.uuid')
    echo "  UUID:   $uuid"
    echo "  State:  $(echo "$order" | jq -r '.state')"

    state="wait"
    for _ in $(seq 1 20); do
      info=$(upbit orders retrieve --uuid "$uuid")
      state=$(echo "$info" | jq -r '.state')
      [ "$state" = "done" ] || [ "$state" = "cancel" ] && break
      sleep 0.5
    done

    if [ "$state" = "done" ]; then
      spent=$(echo "$info" | jq '[.trades[].funds | tonumber] | add // 0')
      volume=$(echo "$info" | jq '[.trades[].volume | tonumber] | add // 0')
      avg=$(echo "scale=0; $spent / $volume" | bc 2>/dev/null || echo "0")
      printf "  Filled amount:  %.0f KRW\n" "$spent"
      printf "  Filled volume:  %.8f\n" "$volume"
      printf "  Filled price:   %.0f KRW\n" "$avg"
      total_spent=$(echo "$total_spent + $spent" | bc)
      total_volume=$(echo "$total_volume + $volume" | bc)
      success_count=$((success_count + 1))
    else
      echo "  Final state: $state"
    fi

    [ "$r" -lt "$TOTAL_ROUNDS" ] && echo "  Waiting ${INTERVAL}s before next buy..." && sleep "$INTERVAL"
  done
}

# ---------------------------------------------------------------------------
# 3. Summary
# ---------------------------------------------------------------------------

section_show_summary() {
  divider "3. Summary"

  if [ "$DRY_RUN" != "false" ]; then
    echo "  No results (Dry run mode)"
    return 0
  fi

  echo "  Successes:        ${success_count}/${TOTAL_ROUNDS}"
  if (( $(echo "$total_volume > 0" | bc -l) )); then
    avg=$(echo "scale=0; $total_spent / $total_volume" | bc)
    printf "  Total spent:      %.0f KRW\n" "$total_spent"
    printf "  Total volume:     %.8f\n" "$total_volume"
    printf "  Avg buy price:    %.0f KRW\n" "$avg"
  else
    echo "  No successful buys."
  fi
}

# ---------------------------------------------------------------------------
# Run
# ---------------------------------------------------------------------------

section_setup || true
section_execute_dca || true
section_show_summary || true

echo ""
printf '%0.s=' {1..60}; echo ""
echo "  DCA scenario completed"
printf '%0.s=' {1..60}; echo ""
echo ""
echo "  * This example is for educational purposes only."
echo "    Use at your own risk for actual trading."
