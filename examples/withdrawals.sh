#!/usr/bin/env bash
# Withdrawal Management Scenario.
#
# Demonstrates withdrawal management features using the Upbit Exchange API.
#   - List allowed withdrawal addresses
#   - Check deposit/withdrawal service status
#   - Fetch withdrawal availability
#   - Fetch withdrawal history (list / single)
#   - Digital asset withdrawal / KRW fiat withdrawal
#
# Dry run (default):
#   Read-only lookups of withdrawal-related information.
#
# Live run (DRY_RUN=false):
#   Edit this file to set the address/amount before running a real withdrawal.
#   WARNING: Withdrawals cannot be reversed. Mis-typed addresses will lose funds.
#
# Usage:
#   UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/withdrawals.sh
#   DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/withdrawals.sh

set -euo pipefail

command -v jq >/dev/null 2>&1 || { echo "jq is required." >&2; exit 1; }

[ -n "${UPBIT_ACCESS_KEY+x}" ] && export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
[ -n "${UPBIT_SECRET_KEY+x}" ] && export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

DRY_RUN="${DRY_RUN:-true}"

divider() {
  echo ""
  printf '%0.s=' {1..60}; echo ""
  echo "  $1"
  printf '%0.s=' {1..60}; echo ""
}

# ---------------------------------------------------------------------------
# 1. List allowed withdrawal addresses
# ---------------------------------------------------------------------------

section_list_withdrawal_addresses() {
  divider "1. Allowed Withdrawal Addresses"

  addresses=$(upbit withdraws list-coin-addresses)
  count=$(echo "$addresses" | jq 'length')

  if [ "$count" = "0" ]; then
    echo "  No withdrawal addresses registered."
    echo "  Register allowed withdrawal addresses in the Upbit web UI first."
    return 0
  fi

  printf "  %-8s %-8s %-12s %s\n" "Currency" "Network" "Exchange" "Address"
  echo "  ------------------------------------------------------------"
  echo "$addresses" | jq -r '.[:10][] | [.currency, .net_type, (.exchange_name // "-"), .withdraw_address] | @tsv' | \
    awk -F'\t' '{ printf "  %-8s %-8s %-12s %s\n", $1, $2, $3, $4 }'
}

# ---------------------------------------------------------------------------
# 2. Wallet service status
# ---------------------------------------------------------------------------

section_check_wallet_status() {
  divider "2. Deposit/Withdrawal Service Status"

  wallets=$(upbit wallet-status list)
  total=$(echo "$wallets" | jq 'length')
  echo "  Total wallets: $total"
  echo ""
  printf "  %-8s %-10s %-12s %s\n" "Currency" "Network" "Wallet" "Block"
  echo "  ------------------------------------------------"
  echo "$wallets" | jq -r '.[] | select(.currency == "BTC" or .currency == "ETH" or .currency == "USDT") |
    [.currency, .net_type, .wallet_state, .block_state] | @tsv' | \
    awk -F'\t' '{ printf "  %-8s %-10s %-12s %s\n", $1, $2, $3, $4 }'
}

# ---------------------------------------------------------------------------
# 3. Withdrawal availability
# ---------------------------------------------------------------------------

section_check_withdrawal_chance() {
  divider "3. BTC Withdrawal Availability"

  chance=$(upbit withdraws retrieve-chance --currency "BTC" --net-type "BTC")
  echo "  Balance:          $(echo "$chance" | jq -r '.account.balance')"
  echo "  Withdraw fee:     $(echo "$chance" | jq -r '.currency.withdraw_fee')"
  echo "  Minimum amount:   $(echo "$chance" | jq -r '.withdraw_limit.minimum')"
  fiat=$(echo "$chance" | jq -r '.withdraw_limit.fiat_currency')
  remaining=$(echo "$chance" | jq -r '.withdraw_limit.remaining_daily_fiat')
  echo "  Remaining today:  $remaining $fiat"
}

# ---------------------------------------------------------------------------
# 4. Withdrawal history
# ---------------------------------------------------------------------------

section_list_withdrawals() {
  divider "4. Withdrawal History"

  withdrawals=$(upbit withdraws list --limit 5 --max-items 5 | jq -s '.')
  count=$(echo "$withdrawals" | jq 'length')

  if [ "$count" = "0" ]; then
    echo "  No withdrawal history."
    return 0
  fi

  printf "  %-10s %-6s %14s %-12s\n" "UUID" "Currency" "Amount" "State"
  echo "  --------------------------------------------------"
  echo "$withdrawals" | jq -r '.[] | [.uuid[:8], .currency, .amount, .state] | @tsv' | \
    awk -F'\t' '{ printf "  %s... %-6s %14s %-12s\n", $1, $2, $3, $4 }'
}

# ---------------------------------------------------------------------------
# 5. Single withdrawal lookup
# ---------------------------------------------------------------------------

section_retrieve_withdrawal() {
  divider "5. Single Withdrawal Lookup"

  if withdrawal=$(upbit withdraws retrieve 2>/dev/null); then
    echo "  UUID:     $(echo "$withdrawal" | jq -r '.uuid')"
    echo "  Currency: $(echo "$withdrawal" | jq -r '.currency')"
    echo "  Amount:   $(echo "$withdrawal" | jq -r '.amount')"
    echo "  State:    $(echo "$withdrawal" | jq -r '.state')"
    txid=$(echo "$withdrawal" | jq -r '.txid // empty')
    [ -n "$txid" ] && echo "  TxID:     $txid"
  else
    echo "  Failed (no withdrawal history may exist)."
  fi
}

# ---------------------------------------------------------------------------
# 6. Digital asset withdrawal
# ---------------------------------------------------------------------------

section_withdraw_coin() {
  divider "6. Digital Asset Withdrawal"

  CURRENCY="USDT"
  NET_TYPE="TRX"
  AMOUNT="13.241"
  ADDRESS="<your_withdrawal_address>"

  echo "  Currency: $CURRENCY"
  echo "  Network:  $NET_TYPE"
  echo "  Amount:   $AMOUNT"
  echo "  Address:  $ADDRESS"

  if [ "$DRY_RUN" != "false" ]; then
    echo ""
    echo "  [Skipped] Dry run mode -- withdrawal omitted."
    return 0
  fi

  if [[ "$ADDRESS" == "<"* ]]; then
    echo ""
    echo "  [Skipped] Withdrawal address not configured. Edit this file directly."
    return 0
  fi

  result=$(upbit withdraws create-withdrawal \
    --currency "$CURRENCY" --net-type "$NET_TYPE" \
    --amount "$AMOUNT" --address "$ADDRESS" --transaction-type "default")
  echo ""
  echo "  UUID:   $(echo "$result" | jq -r '.uuid')"
  echo "  State:  $(echo "$result" | jq -r '.state')"
}

# ---------------------------------------------------------------------------
# 7. KRW fiat withdrawal
# ---------------------------------------------------------------------------

section_withdraw_krw() {
  divider "7. KRW Fiat Withdrawal"

  AMOUNT="10000"
  printf "  Amount: %.0f KRW\n" "$AMOUNT"

  if [ "$DRY_RUN" != "false" ]; then
    echo ""
    echo "  [Skipped] Dry run mode -- withdrawal omitted."
    return 0
  fi

  result=$(upbit withdraws create-krw-withdrawal --amount "$AMOUNT" --two-factor-type "kakao")
  echo ""
  echo "  UUID:   $(echo "$result" | jq -r '.uuid')"
  echo "  State:  $(echo "$result" | jq -r '.state')"
}

# ---------------------------------------------------------------------------
# Run
# ---------------------------------------------------------------------------

section_list_withdrawal_addresses || true
section_check_wallet_status || true
section_check_withdrawal_chance || true
section_list_withdrawals || true
section_retrieve_withdrawal || true
section_withdraw_coin || true
section_withdraw_krw || true

echo ""
printf '%0.s=' {1..60}; echo ""
echo "  Withdrawals scenario completed"
printf '%0.s=' {1..60}; echo ""
echo ""
echo "  * This example is for educational purposes only."
echo "    Use at your own risk for actual trading."
