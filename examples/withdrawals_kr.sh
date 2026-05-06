#!/usr/bin/env bash
# 출금 관리 시나리오.
#
# 업비트 Exchange API를 사용하여 출금 관리 기능을 확인하는 예제입니다.
#   - 출금 허용 주소 목록 조회
#   - 입출금 서비스 상태 확인
#   - 출금 가능 정보 조회
#   - 출금 내역 조회 (목록/개별)
#   - 디지털 자산 출금 / 원화 출금
#
# Dry run (기본):
#   출금 관련 정보 조회만 수행합니다.
#
# 실제 실행 (DRY_RUN=false):
#   실제 출금은 코드에서 주소/금액을 수정한 뒤 실행하세요.
#   주의: 출금은 비가역적입니다. 잘못된 주소로 출금하면 자산을 잃을 수 있습니다.
#
# 실행:
#   UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/withdrawals.sh
#   DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/withdrawals.sh

set -euo pipefail

command -v jq >/dev/null 2>&1 || { echo "jq가 필요합니다." >&2; exit 1; }

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
# 1. 출금 허용 주소 목록 조회
# ---------------------------------------------------------------------------

section_list_withdrawal_addresses() {
  divider "1. 출금 허용 주소 목록"

  addresses=$(upbit withdraws list-coin-addresses)
  count=$(echo "$addresses" | jq 'length')

  if [ "$count" = "0" ]; then
    echo "  등록된 출금 주소가 없습니다."
    echo "  업비트 웹에서 출금 허용 주소를 먼저 등록해주세요."
    return 0
  fi

  printf "  %-8s %-8s %-12s %s\n" "통화" "네트워크" "거래소" "주소"
  echo "  ------------------------------------------------------------"
  echo "$addresses" | jq -r '.[:10][] | [.currency, .net_type, (.exchange_name // "-"), .withdraw_address] | @tsv' | \
    awk -F'\t' '{ printf "  %-8s %-8s %-12s %s\n", $1, $2, $3, $4 }'
}

# ---------------------------------------------------------------------------
# 2. 입출금 서비스 상태 확인
# ---------------------------------------------------------------------------

section_check_wallet_status() {
  divider "2. 입출금 서비스 상태"

  wallets=$(upbit wallet-status list)
  total=$(echo "$wallets" | jq 'length')
  echo "  전체 지갑 수: $total"
  echo ""
  printf "  %-8s %-10s %-12s %s\n" "통화" "네트워크" "상태" "블록 상태"
  echo "  ------------------------------------------------"
  echo "$wallets" | jq -r '.[] | select(.currency == "BTC" or .currency == "ETH" or .currency == "USDT") |
    [.currency, .net_type, .wallet_state, .block_state] | @tsv' | \
    awk -F'\t' '{ printf "  %-8s %-10s %-12s %s\n", $1, $2, $3, $4 }'
}

# ---------------------------------------------------------------------------
# 3. 출금 가능 정보 조회
# ---------------------------------------------------------------------------

section_check_withdrawal_chance() {
  divider "3. BTC 출금 가능 정보"

  chance=$(upbit withdraws retrieve-chance --currency "BTC" --net-type "BTC")
  echo "  보유 수량:     $(echo "$chance" | jq -r '.account.balance')"
  echo "  출금 수수료:   $(echo "$chance" | jq -r '.currency.withdraw_fee')"
  echo "  최소 출금금액: $(echo "$chance" | jq -r '.withdraw_limit.minimum')"
  fiat=$(echo "$chance" | jq -r '.withdraw_limit.fiat_currency')
  remaining=$(echo "$chance" | jq -r '.withdraw_limit.remaining_daily_fiat')
  echo "  1일 잔여한도:  $remaining $fiat"
}

# ---------------------------------------------------------------------------
# 4. 출금 내역 조회
# ---------------------------------------------------------------------------

section_list_withdrawals() {
  divider "4. 출금 내역 조회"

  withdrawals=$(upbit withdraws list --limit 5 --max-items 5 | jq -s '.')
  count=$(echo "$withdrawals" | jq 'length')

  if [ "$count" = "0" ]; then
    echo "  출금 내역이 없습니다."
    return 0
  fi

  printf "  %-10s %-6s %14s %-12s\n" "UUID" "통화" "금액" "상태"
  echo "  --------------------------------------------------"
  echo "$withdrawals" | jq -r '.[] | [.uuid[:8], .currency, .amount, .state] | @tsv' | \
    awk -F'\t' '{ printf "  %s... %-6s %14s %-12s\n", $1, $2, $3, $4 }'
}

# ---------------------------------------------------------------------------
# 5. 개별 출금 조회
# ---------------------------------------------------------------------------

section_retrieve_withdrawal() {
  divider "5. 개별 출금 조회"

  if withdrawal=$(upbit withdraws retrieve 2>/dev/null); then
    echo "  UUID:   $(echo "$withdrawal" | jq -r '.uuid')"
    echo "  통화:   $(echo "$withdrawal" | jq -r '.currency')"
    echo "  금액:   $(echo "$withdrawal" | jq -r '.amount')"
    echo "  상태:   $(echo "$withdrawal" | jq -r '.state')"
    txid=$(echo "$withdrawal" | jq -r '.txid // empty')
    [ -n "$txid" ] && echo "  TxID:   $txid"
  else
    echo "  출금 조회 실패 (출금 내역이 없을 수 있음)"
  fi
}

# ---------------------------------------------------------------------------
# 6. 디지털 자산 출금
# ---------------------------------------------------------------------------

section_withdraw_coin() {
  divider "6. 디지털 자산 출금"

  CURRENCY="USDT"
  NET_TYPE="TRX"
  AMOUNT="13.241"
  ADDRESS="<your_withdrawal_address>"

  echo "  통화:     $CURRENCY"
  echo "  네트워크: $NET_TYPE"
  echo "  금액:     $AMOUNT"
  echo "  주소:     $ADDRESS"

  if [ "$DRY_RUN" != "false" ]; then
    echo ""
    echo "  [건너뜀] Dry run 모드이므로 실제 출금을 생략합니다."
    return 0
  fi

  if [[ "$ADDRESS" == "<"* ]]; then
    echo ""
    echo "  [건너뜀] 출금 주소가 설정되지 않았습니다. 코드에서 직접 수정하세요."
    return 0
  fi

  result=$(upbit withdraws create-withdrawal \
    --currency "$CURRENCY" --net-type "$NET_TYPE" \
    --amount "$AMOUNT" --address "$ADDRESS" --transaction-type "default")
  echo ""
  echo "  출금 UUID: $(echo "$result" | jq -r '.uuid')"
  echo "  상태:      $(echo "$result" | jq -r '.state')"
}

# ---------------------------------------------------------------------------
# 7. 원화(KRW) 출금
# ---------------------------------------------------------------------------

section_withdraw_krw() {
  divider "7. 원화(KRW) 출금"

  AMOUNT="10000"
  printf "  출금 금액: %.0f KRW\n" "$AMOUNT"

  if [ "$DRY_RUN" != "false" ]; then
    echo ""
    echo "  [건너뜀] Dry run 모드이므로 실제 출금을 생략합니다."
    return 0
  fi

  result=$(upbit withdraws create-krw-withdrawal --amount "$AMOUNT" --two-factor-type "kakao")
  echo ""
  echo "  출금 UUID: $(echo "$result" | jq -r '.uuid')"
  echo "  상태:      $(echo "$result" | jq -r '.state')"
}

# ---------------------------------------------------------------------------
# 실행
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
echo "  출금 관리 시나리오 완료"
printf '%0.s=' {1..60}; echo ""
echo ""
echo "  ※ 이 예제는 교육 목적으로 작성되었습니다."
echo "    실제 투자에 사용할 경우 발생하는 손실에 대해 책임지지 않습니다."
