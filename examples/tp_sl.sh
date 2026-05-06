#!/usr/bin/env bash
# Take-Profit / Stop-Loss Auto Sell Scenario.
#
# Polls the current price via the Upbit Exchange API and executes a market sell
# when the target (take-profit) or stop-loss price is reached.
#   - Fetch current price and calculate TP/SL prices
#   - Price monitoring via REST polling
#   - Execute market sell
#
# The Upbit API does not support reserved (stop/limit) orders, so this example
# implements it client-side via polling.
#
# Dry run (default):
#   Uses the current market price as the reference. No orders are placed.
#
# Live run (DRY_RUN=false):
#   Uses the average buy price as the reference and places a real sell order.
#   You must hold the asset.
#
# Usage:
#   UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/tp_sl.sh
#   DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/tp_sl.sh

set -euo pipefail

command -v jq >/dev/null 2>&1 || { echo "jq is required." >&2; exit 1; }
command -v bc >/dev/null 2>&1 || { echo "bc is required." >&2; exit 1; }

[ -n "${UPBIT_ACCESS_KEY+x}" ] && export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
[ -n "${UPBIT_SECRET_KEY+x}" ] && export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

MARKET="KRW-BTC"
SELL_VOLUME="0.0001"
TP_PERCENT="3"    # Take-profit threshold: +3%
SL_PERCENT="2"    # Stop-loss threshold: -2%
POLL_INTERVAL=1   # Poll interval (seconds)
MAX_POLLS=10      # Maximum number of polls
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
# 1. TP/SL configuration
# ---------------------------------------------------------------------------

section_setup() {
  local mode_label
  [ "$DRY_RUN" = "false" ] && mode_label="LIVE RUN" || mode_label="DRY RUN"
  divider "1. TP/SL Configuration [$mode_label]"

  local current
  current=$(current_price)

  if [ "$DRY_RUN" = "false" ]; then
    coin=$(echo "$MARKET" | cut -d'-' -f2)
    accounts=$(upbit accounts list)
    target=$(echo "$accounts" | jq --arg c "$coin" '.[] | select(.currency == $c)')

    if [ -z "$target" ]; then
      echo "  [Error] No $coin holdings." >&2; exit 1
    fi

    balance=$(echo "$target" | jq -r '.balance')
    if (( $(echo "$balance < $SELL_VOLUME" | bc -l) )); then
      echo "  [Error] $coin balance ($balance) is less than sell volume ($SELL_VOLUME)." >&2; exit 1
    fi

    BASE_PRICE=$(echo "$target" | jq -r '.avg_buy_price')
    echo "  Balance:       $coin $balance"
    printf "  Base price:    %.0f KRW (avg buy price)\n" "$BASE_PRICE"
    printf "  Current:       %.0f KRW\n" "$current"
  else
    BASE_PRICE="$current"
    printf "  Base price:    %.0f KRW (current price)\n" "$BASE_PRICE"
  fi

  TP_PRICE=$(echo "scale=0; $BASE_PRICE * (1 + $TP_PERCENT / 100) / 1" | bc)
  SL_PRICE=$(echo "scale=0; $BASE_PRICE * (1 - $SL_PERCENT / 100) / 1" | bc)

  echo "  Market:        $MARKET"
  echo "  Sell volume:   $SELL_VOLUME"
  printf "  TP price:      %.0f KRW (+%s%%)\n" "$TP_PRICE" "$TP_PERCENT"
  printf "  SL price:      %.0f KRW (-%s%%)\n" "$SL_PRICE" "$SL_PERCENT"
  echo "  Poll interval: ${POLL_INTERVAL}s x ${MAX_POLLS} max"
  [ "$DRY_RUN" != "false" ] && echo "" && echo "  * Dry run mode: no orders are placed."
}

# ---------------------------------------------------------------------------
# 2. Price monitoring
# ---------------------------------------------------------------------------

section_monitor_price() {
  divider "2. Price Monitoring"

  TRIGGER="timeout"

  for poll in $(seq 1 "$MAX_POLLS"); do
    current=$(current_price)
    printf "  [%3d/%d] Price: %.0f KRW" "$poll" "$MAX_POLLS" "$current"

    if (( $(echo "$current >= $TP_PRICE" | bc -l) )); then
      echo "  -> TP reached!"
      TRIGGER="tp"; break
    fi
    if (( $(echo "$current <= $SL_PRICE" | bc -l) )); then
      echo "  -> SL reached!"
      TRIGGER="sl"; break
    fi

    echo ""
    [ "$poll" -lt "$MAX_POLLS" ] && sleep "$POLL_INTERVAL"
  done

  if [ "$TRIGGER" = "timeout" ]; then
    echo ""
    echo "  Max polls (${MAX_POLLS}) reached."
  fi
}

# ---------------------------------------------------------------------------
# 3. Execute market sell
# ---------------------------------------------------------------------------

section_execute_sell() {
  local label
  [ "$TRIGGER" = "tp" ] && label="TP" || label="SL"
  divider "3. Execute Sell (Trigger: $label)"

  if [ "$DRY_RUN" != "false" ]; then
    echo "  [Skipped] Dry run mode -- actual order omitted."
    return 0
  fi

  order=$(upbit orders create \
    --market "$MARKET" --side "ask" --ord-type "market" --volume "$SELL_VOLUME")
  uuid=$(echo "$order" | jq -r '.uuid')
  echo "  Order created"
  echo "  UUID:   $uuid"
  echo "  State:  $(echo "$order" | jq -r '.state')"

  state="wait"
  for _ in $(seq 1 20); do
    info=$(upbit orders retrieve --uuid "$uuid")
    state=$(echo "$info" | jq -r '.state')
    [ "$state" = "done" ] || [ "$state" = "cancel" ] && break
    sleep 0.5
  done

  echo "  Final state:     $state"
  echo "  Filled volume:   $(echo "$info" | jq -r '.executed_volume')"
  if [ "$state" = "done" ]; then
    total=$(echo "$info" | jq '[.trades[].funds | tonumber] | add // 0')
    printf "  Filled amount:   %.0f KRW\n" "$total"
  fi
}

# ---------------------------------------------------------------------------
# Run
# ---------------------------------------------------------------------------

section_setup || true
section_monitor_price || true

if [ "$TRIGGER" = "timeout" ]; then
  echo ""
  echo "  Skipping sell due to timeout."
else
  section_execute_sell
fi

echo ""
printf '%0.s=' {1..60}; echo ""
echo "  TP/SL scenario completed"
printf '%0.s=' {1..60}; echo ""
echo ""
echo "  * This example is for educational purposes only."
echo "    Use at your own risk for actual trading."
