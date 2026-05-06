#!/usr/bin/env bash
# 입금 관리 시나리오.
#
# 업비트 Exchange API를 사용하여 입금 관리 기능을 확인하는 예제입니다.
#   - 입금 주소 목록 조회
#   - 입금 주소 조회 / 생성
#   - 입금 가능 정보 조회
#   - 입금 내역 조회 (목록/개별)
#   - 트래블룰 지원 거래소 목록 조회
#   - 트래블룰 검증 (입금 UUID 기준)
#
# Dry run (기본):
#   입금 관련 정보 조회만 수행합니다.
#
# 실제 실행 (DRY_RUN=false):
#   입금 주소 생성 및 트래블룰 검증을 실행합니다.
#
# 실행:
#   UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/deposits.sh
#   DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/deposits.sh

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
# 1. 입금 주소 목록 조회
# ---------------------------------------------------------------------------

section_list_deposit_addresses() {
  divider "1. 입금 주소 목록 조회"

  addresses=$(upbit deposits list-coin-addresses)
  count=$(echo "$addresses" | jq 'length')

  if [ "$count" = "0" ]; then
    echo "  등록된 입금 주소가 없습니다."
    return 0
  fi

  printf "  %-8s %-8s %s\n" "통화" "네트워크" "주소"
  echo "  --------------------------------------------------"
  echo "$addresses" | jq -r '.[:10][] | [.currency, .net_type, (.deposit_address // "(미생성)")] | @tsv' | \
    awk -F'\t' '{ printf "  %-8s %-8s %s\n", $1, $2, $3 }'
}

# ---------------------------------------------------------------------------
# 2. 특정 코인 입금 주소 조회
# ---------------------------------------------------------------------------

section_retrieve_deposit_address() {
  divider "2. BTC 입금 주소 조회"

  if address=$(upbit deposits retrieve-coin-address --currency "BTC" --net-type "BTC" 2>/dev/null); then
    deposit_addr=$(echo "$address" | jq -r '.deposit_address // empty')
    if [ -n "$deposit_addr" ]; then
      echo "  통화:     $(echo "$address" | jq -r '.currency')"
      echo "  네트워크: $(echo "$address" | jq -r '.net_type')"
      echo "  주소:     $deposit_addr"
      secondary=$(echo "$address" | jq -r '.secondary_address // empty')
      [ -n "$secondary" ] && echo "  2차 주소: $secondary"
    else
      echo "  BTC 입금 주소가 아직 생성되지 않았습니다."
    fi
  else
    echo "  입금 주소 조회 실패"
  fi
}

# ---------------------------------------------------------------------------
# 3. 입금 주소 생성
# ---------------------------------------------------------------------------

section_create_deposit_address() {
  divider "3. USDT(TRX) 입금 주소 생성"

  if [ "$DRY_RUN" != "false" ]; then
    echo "  [건너뜀] Dry run 모드이므로 실제 주소 생성을 생략합니다."
    echo "  생성 대상: USDT (TRX)"
    return 0
  fi

  if result=$(upbit deposits create-coin-address --currency "USDT" --net-type "TRX" 2>/dev/null); then
    deposit_addr=$(echo "$result" | jq -r '.deposit_address // empty')
    if [ -n "$deposit_addr" ]; then
      echo "  통화:     $(echo "$result" | jq -r '.currency')"
      echo "  네트워크: $(echo "$result" | jq -r '.net_type')"
      echo "  주소:     $deposit_addr"
    else
      echo "  생성 요청 상태: $(echo "$result" | jq -r '.success')"
      echo "  메시지:         $(echo "$result" | jq -r '.message')"
    fi
  else
    echo "  입금 주소 생성 실패"
  fi
}

# ---------------------------------------------------------------------------
# 4. 입금 가능 정보 조회
# ---------------------------------------------------------------------------

section_retrieve_deposit_chance() {
  divider "4. BTC 입금 가능 정보"

  chance=$(upbit deposits retrieve-chance --currency "BTC" --net-type "BTC")
  possible=$(echo "$chance" | jq -r '.is_deposit_possible')
  echo "  입금 가능:      $([ "$possible" = "true" ] && echo "예" || echo "아니오")"
  echo "  최소 입금금액:  $(echo "$chance" | jq -r '.minimum_deposit_amount')"
  echo "  필요 확인 수:   $(echo "$chance" | jq -r '.minimum_deposit_confirmations')"
  if [ "$possible" = "false" ]; then
    echo "  불가 사유:      $(echo "$chance" | jq -r '.deposit_impossible_reason')"
  fi
}

# ---------------------------------------------------------------------------
# 5. 입금 내역 조회
# ---------------------------------------------------------------------------

section_list_deposits() {
  divider "5. 입금 내역 조회"

  deposits=$(upbit deposits list --limit 5 --max-items 5 | jq -s '.')
  count=$(echo "$deposits" | jq 'length')

  if [ "$count" = "0" ]; then
    echo "  입금 내역이 없습니다."
    return 0
  fi

  printf "  %-10s %-6s %14s %-24s\n" "UUID" "통화" "금액" "상태"
  echo "  ------------------------------------------------------------"
  echo "$deposits" | jq -r '.[] | [.uuid[:8], .currency, .amount, .state] | @tsv' | \
    awk -F'\t' '{ printf "  %s... %-6s %14s %-24s\n", $1, $2, $3, $4 }'
}

# ---------------------------------------------------------------------------
# 6. 개별 입금 조회
# ---------------------------------------------------------------------------

section_retrieve_deposit() {
  divider "6. 개별 입금 조회"

  if deposit=$(upbit deposits retrieve 2>/dev/null); then
    uuid=$(echo "$deposit" | jq -r '.uuid')
    echo "  UUID:   $uuid"
    echo "  통화:   $(echo "$deposit" | jq -r '.currency')"
    echo "  금액:   $(echo "$deposit" | jq -r '.amount')"
    echo "  상태:   $(echo "$deposit" | jq -r '.state')"
    txid=$(echo "$deposit" | jq -r '.txid // empty')
    [ -n "$txid" ] && echo "  TxID:   $txid"

    # UUID로 재조회
    echo ""
    same=$(upbit deposits retrieve --uuid "$uuid")
    echo "  UUID 재조회 상태: $(echo "$same" | jq -r '.state')"
  else
    echo "  입금 조회 실패 (입금 내역이 없을 수 있음)"
  fi
}

# ---------------------------------------------------------------------------
# 7. 트래블룰 지원 거래소 목록 조회
# ---------------------------------------------------------------------------

section_list_travel_rule_vasps() {
  divider "7. 트래블룰 지원 거래소 목록"

  vasps=$(upbit travel-rule list-vasps)
  total=$(echo "$vasps" | jq 'length')
  echo "  지원 거래소 수: $total"
  echo ""
  printf "  %-20s %s\n" "거래소명" "UUID"
  echo "  --------------------------------------------------"
  echo "$vasps" | jq -r '.[:10][] | [.vasp_name, .vasp_uuid] | @tsv' | \
    awk -F'\t' '{ printf "  %-20s %s\n", $1, $2 }'
  [ "$total" -gt 10 ] && echo "  ... 외 $((total - 10))개"
}

# ---------------------------------------------------------------------------
# 8. 트래블룰 검증
# ---------------------------------------------------------------------------

section_verify_travel_rule() {
  divider "8. 트래블룰 검증"

  DEPOSIT_UUID="<your_deposit_uuid>"
  VASP_UUID="<deposit_vasp_uuid>"

  echo "  Deposit UUID: $DEPOSIT_UUID"
  echo "  거래소 UUID:   $VASP_UUID"

  if [ "$DRY_RUN" != "false" ]; then
    echo ""
    echo "  [건너뜀] Dry run 모드이므로 트래블룰 검증을 생략합니다."
    return 0
  fi

  if [[ "$DEPOSIT_UUID" == "<"* ]] || [[ "$VASP_UUID" == "<"* ]]; then
    echo ""
    echo "  [건너뜀] UUID가 설정되지 않았습니다. 코드에서 직접 수정하세요."
    return 0
  fi

  result=$(upbit travel-rule verify-deposit-by-uuid \
    --deposit-uuid "$DEPOSIT_UUID" --vasp-uuid "$VASP_UUID")
  echo "  입금 UUID:   $(echo "$result" | jq -r '.deposit_uuid')"
  echo "  입금 상태:   $(echo "$result" | jq -r '.deposit_state')"
  echo "  검증 결과:   $(echo "$result" | jq -r '.verification_result')"
}

# ---------------------------------------------------------------------------
# 실행
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
echo "  입금 관리 시나리오 완료"
printf '%0.s=' {1..60}; echo ""
echo ""
echo "  ※ 이 예제는 교육 목적으로 작성되었습니다."
echo "    실제 투자에 사용할 경우 발생하는 손실에 대해 책임지지 않습니다."
