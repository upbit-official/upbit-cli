#!/usr/bin/env bash
# Deposit Management Scenario.
#
# Demonstrates deposit management features using the Upbit Exchange API.
#   - List deposit addresses
#   - Look up / create a deposit address
#   - Check deposit availability
#   - Fetch deposit history (list / single)
#   - List Travel Rule-supported exchanges
#   - Verify Travel Rule (by deposit UUID)
#
# Dry run (default):
#   Read-only lookups of deposit-related information.
#
# Live run (DRY_RUN=false):
#   Performs deposit address generation and Travel Rule verification.
#
# Usage:
#   UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/deposits.sh
#   DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/deposits.sh

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
# 1. List deposit addresses
# ---------------------------------------------------------------------------

section_list_deposit_addresses() {
  divider "1. Deposit Addresses"

  addresses=$(upbit deposits list-coin-addresses)
  count=$(echo "$addresses" | jq 'length')

  if [ "$count" = "0" ]; then
    echo "  No deposit addresses registered."
    return 0
  fi

  printf "  %-8s %-8s %s\n" "Currency" "Network" "Address"
  echo "  --------------------------------------------------"
  echo "$addresses" | jq -r '.[:10][] | [.currency, .net_type, (.deposit_address // "(not created)")] | @tsv' | \
    awk -F'\t' '{ printf "  %-8s %-8s %s\n", $1, $2, $3 }'
}

# ---------------------------------------------------------------------------
# 2. Retrieve deposit address
# ---------------------------------------------------------------------------

section_retrieve_deposit_address() {
  divider "2. Retrieve BTC Deposit Address"

  if address=$(upbit deposits retrieve-coin-address --currency "BTC" --net-type "BTC" 2>/dev/null); then
    deposit_addr=$(echo "$address" | jq -r '.deposit_address // empty')
    if [ -n "$deposit_addr" ]; then
      echo "  Currency: $(echo "$address" | jq -r '.currency')"
      echo "  Network:  $(echo "$address" | jq -r '.net_type')"
      echo "  Address:  $deposit_addr"
      secondary=$(echo "$address" | jq -r '.secondary_address // empty')
      [ -n "$secondary" ] && echo "  Tag:      $secondary"
    else
      echo "  BTC deposit address has not been created yet."
    fi
  else
    echo "  Failed to retrieve deposit address."
  fi
}

# ---------------------------------------------------------------------------
# 3. Create deposit address
# ---------------------------------------------------------------------------

section_create_deposit_address() {
  divider "3. Create USDT(TRX) Deposit Address"

  if [ "$DRY_RUN" != "false" ]; then
    echo "  [Skipped] Dry run mode -- actual address creation omitted."
    echo "  Target: USDT (TRX)"
    return 0
  fi

  if result=$(upbit deposits create-coin-address --currency "USDT" --net-type "TRX" 2>/dev/null); then
    deposit_addr=$(echo "$result" | jq -r '.deposit_address // empty')
    if [ -n "$deposit_addr" ]; then
      echo "  Currency: $(echo "$result" | jq -r '.currency')"
      echo "  Network:  $(echo "$result" | jq -r '.net_type')"
      echo "  Address:  $deposit_addr"
    else
      echo "  Request state: $(echo "$result" | jq -r '.success')"
      echo "  Message:       $(echo "$result" | jq -r '.message')"
    fi
  else
    echo "  Failed to create deposit address."
  fi
}

# ---------------------------------------------------------------------------
# 4. Deposit availability
# ---------------------------------------------------------------------------

section_retrieve_deposit_chance() {
  divider "4. BTC Deposit Availability"

  chance=$(upbit deposits retrieve-chance --currency "BTC" --net-type "BTC")
  possible=$(echo "$chance" | jq -r '.is_deposit_possible')
  echo "  Deposit possible:      $([ "$possible" = "true" ] && echo "yes" || echo "no")"
  echo "  Minimum deposit:       $(echo "$chance" | jq -r '.minimum_deposit_amount')"
  echo "  Required confirmations: $(echo "$chance" | jq -r '.minimum_deposit_confirmations')"
  if [ "$possible" = "false" ]; then
    echo "  Reason:                $(echo "$chance" | jq -r '.deposit_impossible_reason')"
  fi
}

# ---------------------------------------------------------------------------
# 5. Deposit history
# ---------------------------------------------------------------------------

section_list_deposits() {
  divider "5. Deposit History"

  deposits=$(upbit deposits list --limit 5 --max-items 5 | jq -s '.')
  count=$(echo "$deposits" | jq 'length')

  if [ "$count" = "0" ]; then
    echo "  No deposit history."
    return 0
  fi

  printf "  %-10s %-6s %14s %-24s\n" "UUID" "Currency" "Amount" "State"
  echo "  ------------------------------------------------------------"
  echo "$deposits" | jq -r '.[] | [.uuid[:8], .currency, .amount, .state] | @tsv' | \
    awk -F'\t' '{ printf "  %s... %-6s %14s %-24s\n", $1, $2, $3, $4 }'
}

# ---------------------------------------------------------------------------
# 6. Individual deposit lookup
# ---------------------------------------------------------------------------

section_retrieve_deposit() {
  divider "6. Single Deposit Lookup"

  if deposit=$(upbit deposits retrieve 2>/dev/null); then
    uuid=$(echo "$deposit" | jq -r '.uuid')
    echo "  UUID:     $uuid"
    echo "  Currency: $(echo "$deposit" | jq -r '.currency')"
    echo "  Amount:   $(echo "$deposit" | jq -r '.amount')"
    echo "  State:    $(echo "$deposit" | jq -r '.state')"
    txid=$(echo "$deposit" | jq -r '.txid // empty')
    [ -n "$txid" ] && echo "  TxID:     $txid"

    echo ""
    same=$(upbit deposits retrieve --uuid "$uuid")
    echo "  Lookup by UUID: state=$(echo "$same" | jq -r '.state')"
  else
    echo "  Failed (no deposit history may exist)."
  fi
}

# ---------------------------------------------------------------------------
# 7. Travel Rule VASPs
# ---------------------------------------------------------------------------

section_list_travel_rule_vasps() {
  divider "7. Travel Rule Supported Exchanges"

  vasps=$(upbit travel-rule list-vasps)
  total=$(echo "$vasps" | jq 'length')
  echo "  Supported exchanges: $total"
  echo ""
  printf "  %-20s %s\n" "Name" "UUID"
  echo "  --------------------------------------------------"
  echo "$vasps" | jq -r '.[:10][] | [.vasp_name, .vasp_uuid] | @tsv' | \
    awk -F'\t' '{ printf "  %-20s %s\n", $1, $2 }'
  [ "$total" -gt 10 ] && echo "  ... and $((total - 10)) more"
}

# ---------------------------------------------------------------------------
# 8. Travel Rule verification
# ---------------------------------------------------------------------------

section_verify_travel_rule() {
  divider "8. Travel Rule Verification"

  DEPOSIT_UUID="<your_deposit_uuid>"
  VASP_UUID="<deposit_vasp_uuid>"

  echo "  Deposit UUID: $DEPOSIT_UUID"
  echo "  VASP UUID:    $VASP_UUID"

  if [ "$DRY_RUN" != "false" ]; then
    echo ""
    echo "  [Skipped] Dry run mode -- verification omitted."
    return 0
  fi

  if [[ "$DEPOSIT_UUID" == "<"* ]] || [[ "$VASP_UUID" == "<"* ]]; then
    echo ""
    echo "  [Skipped] UUIDs not configured. Edit this file directly."
    return 0
  fi

  result=$(upbit travel-rule verify-deposit-by-uuid \
    --deposit-uuid "$DEPOSIT_UUID" --vasp-uuid "$VASP_UUID")
  echo "  Deposit UUID:         $(echo "$result" | jq -r '.deposit_uuid')"
  echo "  Deposit state:        $(echo "$result" | jq -r '.deposit_state')"
  echo "  Verification result:  $(echo "$result" | jq -r '.verification_result')"
}

# ---------------------------------------------------------------------------
# Run
# ---------------------------------------------------------------------------

section_list_deposit_addresses || true
section_retrieve_deposit_address || true
section_create_deposit_address || true
section_retrieve_deposit_chance || true
section_list_deposits || true
section_retrieve_deposit || true
section_list_travel_rule_vasps || true
section_verify_travel_rule || true

echo ""
printf '%0.s=' {1..60}; echo ""
echo "  Deposits scenario completed"
printf '%0.s=' {1..60}; echo ""
echo ""
echo "  * This example is for educational purposes only."
echo "    Use at your own risk for actual trading."
