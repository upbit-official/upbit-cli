#!/usr/bin/env bash
# Indicators Scenario.
#
# Calculates investment indicators using data from the Upbit Quotation API.
#   - Top 5 pairs by 24-hour cumulative trading volume (KRW markets)
#
# The Quotation API does not require authentication.
#
# Note: RSI calculation is not included due to shell environment limitations.
#       For RSI examples, refer to python/examples/indicators.py.
#
# Usage:
#   bash examples/indicators.sh

set -euo pipefail

command -v jq >/dev/null 2>&1 || { echo "jq is required." >&2; exit 1; }

divider() {
  echo ""
  printf '%0.s=' {1..60}; echo ""
  echo "  $1"
  printf '%0.s=' {1..60}; echo ""
}

# ---------------------------------------------------------------------------
# 1. Top 5 pairs by 24-hour cumulative trading volume (KRW markets)
# ---------------------------------------------------------------------------

section_top5() {
  divider "1. Top 5 by 24h Trading Volume (KRW markets)"

  tickers=$(upbit tickers list-by-quote-currencies --quote-currencies "KRW")
  total=$(echo "$tickers" | jq 'length')
  echo "  KRW markets: $total"
  echo ""
  printf "  %-12s %14s %20s\n" "Market" "Price" "24h Volume (KRW)"
  echo "  ----------------------------------------------------------------"

  echo "$tickers" | jq -r 'sort_by(-.acc_trade_price_24h | tonumber) | .[:5][] |
    [.market, .trade_price, .acc_trade_price_24h] | @tsv' | \
    awk -F'\t' '{
      acc = $3 + 0
      if (acc >= 1000000000000)
        unit = sprintf("~%.1fT", acc / 1000000000000)
      else if (acc >= 1000000000)
        unit = sprintf("~%.1fB", acc / 1000000000)
      else if (acc >= 1000000)
        unit = sprintf("~%.1fM", acc / 1000000)
      else
        unit = sprintf("%.0f", acc)
      printf "  %-12s %14.0f %20.0f  (%s)\n", $1, $2, acc, unit
    }'
}

# ---------------------------------------------------------------------------
# Run
# ---------------------------------------------------------------------------

section_top5

echo ""
printf '%0.s=' {1..60}; echo ""
echo "  Indicators scenario completed"
printf '%0.s=' {1..60}; echo ""
echo ""
echo "  * This example is for educational purposes only."
echo "    Use at your own risk for actual trading."
