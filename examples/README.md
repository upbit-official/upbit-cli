English | [한국어](./README_KR.md)

# Example Code

This directory contains examples that demonstrate the key features of the Upbit CLI SDK in scenario-based flows.

Each example is written based on real trading workflows and runs in **Dry run mode** by default.
In dry run mode, only read operations are performed — write operations such as placing orders or making withdrawals are skipped.

> **Note:** Per-endpoint examples are organized in feature-specific subdirectories within this directory (e.g., [`examples/orders/`](./orders/) for order-related examples). For Korean scenario examples, see [`README_KR.md`](README_KR.md).

---

## Quick Start

If you are new, we recommend going through the examples in this order:

1. [`quotation.sh`](./quotation.sh) — Quickly explore market data, trades, and orderbook without authentication
2. [`indicators.sh`](./indicators.sh) — Find top pairs by 24h trading volume using Quotation data
3. [`orders.sh`](./orders.sh) — Walk through order creation, lookup, and cancellation in Dry run first

### Examples that do not require authentication

The following examples can be run immediately without an API key:
- [`quotation.sh`](./quotation.sh)
- [`indicators.sh`](./indicators.sh)

### Examples that require authentication

The following examples require an Upbit API key:
- [`orders.sh`](./orders.sh)
- [`orders_test.sh`](./orders_test.sh)
- [`deposits.sh`](./deposits.sh)
- [`withdrawals.sh`](./withdrawals.sh)
- [`dca.sh`](./dca.sh)
- [`tp_sl.sh`](./tp_sl.sh)

### Prerequisites

```bash
# Install jq (required for JSON parsing)
# macOS
brew install jq
# Ubuntu/Debian
apt-get install jq

# Issue an Upbit API key and set environment variables
export UPBIT_ACCESS_KEY=<your-access-key>
export UPBIT_SECRET_KEY=<your-secret-key>
```

### Before you run

- **To explore read-only behavior first**, start with `quotation.sh` and `indicators.sh`.
- **When reviewing examples that can place orders or make withdrawals**, verify the behavior in Dry run mode first.
- **When running with `DRY_RUN=false`**, write operations such as orders, withdrawals, and automated trading will actually execute.

---

## Example Overview

| Example | Purpose | Auth | Default Behavior |
|---|---|---|---|
| [`quotation.sh`](./quotation.sh) | Market data / candles / trades / orderbook | Not required | Safe read-only |
| [`indicators.sh`](./indicators.sh) | Top pairs by 24h trading volume | Not required | Safe read-only |
| [`orders.sh`](./orders.sh) | Order creation, lookup, and cancellation flow | Required | Dry run by default; real orders possible |
| [`orders_test.sh`](./orders_test.sh) | Validate all order types via Order Creation Test API | Required | Always safe — no real orders placed |
| [`deposits.sh`](./deposits.sh) | Deposit address and history management | Required | Dry run by default; some operations may have real effect |
| [`withdrawals.sh`](./withdrawals.sh) | Withdrawal info and withdrawal flow | Required | Dry run by default; real withdrawals require caution |
| [`dca.sh`](./dca.sh) | Automated recurring market-buy (DCA) | Required | Dry run by default; real orders possible |
| [`tp_sl.sh`](./tp_sl.sh) | Automated take-profit / stop-loss sell | Required | Dry run by default; real orders possible |

---

## Scenario Examples in Detail

### `quotation.sh` — Market Data Primer

A quick introduction to the Quotation API, usable without authentication.

**Covered features:**
- List tradable pairs (markets) and filter for caution-flagged tickers
- Fetch current price (ticker) — single, multiple, and per market
- Fetch candle (OHLCV) data — 5-minute, daily, weekly, monthly / using the `to` parameter
- Fetch recent trade history — using the `days_ago` parameter
- Fetch orderbook data

```bash
bash examples/quotation.sh
```

---

### `orders.sh` — Order Creation and Management

Covers the full flow from placing a limit buy order to looking it up and cancelling it.

```bash
# Dry run (default) — read-only, no orders placed
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/orders.sh

# Live run — limit buy -> lookup -> cancel
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/orders.sh
```

---

### `orders_test.sh` — Order Type Validation (test-create API)

Validates all order types safely using the Order Creation Test API.
No real orders are placed — no fees are charged.

```bash
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/orders_test.sh
```

---

### `deposits.sh` — Deposit Address and History Management

```bash
# Dry run (default) — read-only
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/deposits.sh

# Live run — includes deposit address generation
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/deposits.sh
```

---

### `withdrawals.sh` — Withdrawal Management

> **Warning:** Withdrawals cannot be reversed. Always verify the address and amount before executing a real withdrawal.

```bash
# Dry run (default) — read-only
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/withdrawals.sh
```

---

### `dca.sh` — DCA Automated Recurring Buy

Implements a Dollar Cost Averaging (DCA) strategy that repeatedly buys a fixed amount at market price.

```bash
# Dry run (default) — fetches current price only, no orders placed
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/dca.sh

# Live run
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/dca.sh
```

---

### `tp_sl.sh` — Take-Profit / Stop-Loss Auto Sell

Polls the current price and executes a market sell when the target (take-profit) or stop-loss price is reached.

```bash
# Dry run (default) — monitoring only, no orders placed
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/tp_sl.sh

# Live run — auto sell based on average buy price
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/tp_sl.sh
```

---

### `indicators.sh` — Investment Indicators

Finds the top pairs by 24-hour cumulative trading volume in KRW markets.

```bash
bash examples/indicators.sh
```

> RSI calculation is not included due to shell environment limitations. For RSI examples, refer to [`python/examples/indicators.py`](https://github.com/upbit-official/upbit-sdk-python/blob/main/examples/indicators.py).
