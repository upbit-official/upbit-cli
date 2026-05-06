#!/usr/bin/env bash
# Quotation Scenario.
#
# Demonstrates how to query market data using the Upbit Quotation API.
#   - List available markets (KRW pairs)
#   - Get ticker (current price) — single/multiple/by quote
#   - Get candles (OHLCV) — minutes/days/weeks/months, to parameter
#   - Get recent trades (including days_ago)
#   - Get orderbook
#
# The Quotation API does not require authentication.
#
# Usage:
#     bash examples/quotation.sh

set -euo pipefail

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------
divider() {
  echo ""
  echo "============================================================"
  echo "  $1"
  echo "============================================================"
}

# ---------------------------------------------------------------------------
# 1. List Markets
# ---------------------------------------------------------------------------
list_markets() {
  divider "1. List Markets"

  echo "  [All markets]"
  upbit trading-pairs list --is-details=false | jq -r '
    "  Total markets: \(length)\n",
    "  \("Market" | . + " " * (12 - length))\("Korean" | . + " " * (12 - length))English",
    "  --------------------------------------------------",
    (.[:5][] | "  \(.market | . + " " * (12 - length))\(.korean_name | . + " " * (12 - length))\(.english_name)")
  '

  echo ""
  echo "  [Caution markets]"
  upbit trading-pairs list --is-details=true | jq -r '
    [.[] | select(.market_warning == "CAUTION")] |
    "  Caution markets: \(length)"
  '
}

# ---------------------------------------------------------------------------
# 2. Get Ticker
# ---------------------------------------------------------------------------
get_ticker() {
  divider "2. Get Ticker (Current Price)"

  echo "  [Single market]"
  upbit tickers list-by-trading-pairs --markets "KRW-BTC" | jq -r '
    .[0] |
    "  Market:        \(.market)",
    "  Price:         \(.trade_price) KRW",
    "  24h Volume:    \(.acc_trade_price_24h) KRW",
    "  Change Rate:   \(.signed_change_rate)"
  '

  echo ""
  echo "  [Multiple markets]"
  upbit tickers list-by-trading-pairs --markets "KRW-BTC,KRW-ETH" | jq -r '
    .[] | "  \(.market | . + " " * (12 - length))Price: \(.trade_price)"
  '

  echo ""
  echo "  [All KRW tickers]"
  upbit tickers list-by-quote-currencies --quote-currencies "KRW" | jq -r '
    "  KRW market tickers: \(length)"
  '
}

# ---------------------------------------------------------------------------
# 3. Get Candles
# ---------------------------------------------------------------------------
get_candles() {
  divider "3. Get Candles"

  echo "  [5-min candles]"
  upbit candles list-minutes --unit 5 --market "KRW-BTC" --count 3 | jq -r '
    if length == 0 then "  No candle data available."
    else
      (.[] | "  \(.candle_date_time_kst)  Open: \(.opening_price)  Close: \(.trade_price)  Vol: \(.candle_acc_trade_volume)")
    end
  '

  echo ""
  echo "  [to parameter — 1-min candle before 2025-01-01]"
  upbit candles list-minutes --unit 1 --market "KRW-BTC" --to "2025-01-01T00:00:00Z" --count 1 | jq -r '
    if length == 0 then "  No data for specified time range."
    else .[0] | "  \(.candle_date_time_kst)  Close: \(.trade_price)"
    end
  '

  echo ""
  echo "  [Daily candles]"
  upbit candles list-days --market "KRW-BTC" --count 3 | jq -r '
    if length == 0 then "  No daily candle data available."
    else .[] | "  \(.candle_date_time_kst[:10])  Close: \(.trade_price)"
    end
  '

  echo ""
  echo "  [Weekly candles]"
  upbit candles list-weeks --market "KRW-BTC" --count 3 | jq -r '
    if length == 0 then "  No weekly candle data available."
    else .[] | "  \(.candle_date_time_kst[:10])  Close: \(.trade_price)"
    end
  '

  echo ""
  echo "  [Monthly candles]"
  upbit candles list-months --market "KRW-BTC" --count 3 | jq -r '
    if length == 0 then "  No monthly candle data available."
    else .[] | "  \(.candle_date_time_kst[:10])  Close: \(.trade_price)"
    end
  '
}

# ---------------------------------------------------------------------------
# 4. Get Recent Trades
# ---------------------------------------------------------------------------
get_trades() {
  divider "4. Recent Trades"

  upbit trades list --market "KRW-BTC" --count 5 | jq -r '
    .[] |
    (if .ask_bid == "BID" then "BUY" else "SELL" end) as $side |
    "  \(.trade_date_utc) \(.trade_time_utc)  \($side)  \(.trade_price)  \(.trade_volume)"
  '

  echo ""
  echo "  [Trades from 1 day ago]"
  upbit trades list --market "KRW-BTC" --days-ago 1 --count 3 | jq -r '
    if length == 0 then "  No trades found for 1 day ago."
    else .[] |
      (if .ask_bid == "BID" then "BUY" else "SELL" end) as $side |
      "  \(.trade_date_utc) \(.trade_time_utc)  \($side)  \(.trade_price)"
    end
  '
}

# ---------------------------------------------------------------------------
# 5. Get Orderbook
# ---------------------------------------------------------------------------
get_orderbook() {
  divider "5. Orderbook"

  upbit orderbooks list --markets "KRW-BTC" | jq -r '
    .[0] |
    "  Total Ask Size: \(.total_ask_size)",
    "  Total Bid Size: \(.total_bid_size)",
    "",
    (.orderbook_units[:5][] |
      "  Ask: \(.ask_price) (\(.ask_size))   Bid: \(.bid_price) (\(.bid_size))"
    )
  '
}

# ---------------------------------------------------------------------------
# Run
# ---------------------------------------------------------------------------
list_markets
get_ticker
get_candles
get_trades
get_orderbook

echo ""
echo "============================================================"
echo "  Quotation scenario completed"
echo "============================================================"
echo ""
echo "  * This example is for educational purposes only."
echo "    Use at your own risk for actual trading."
