#!/usr/bin/env bash
# DCA (정기 매수) 자동 시장가 매수 시나리오.
#
# 업비트 Exchange API를 사용하여 일정 금액을 주기적으로 시장가 매수하는 예제입니다.
#   - 현재가 조회 및 매수 설정 확인
#   - 지정 횟수만큼 시장가 매수 반복
#   - 매수 결과 요약 (총 매수금액, 총 매수수량, 평균 단가)
#
# Dry run (기본):
#   현재가 조회만 수행하며, 주문을 실행하지 않습니다.
#
# 실제 실행 (DRY_RUN=false):
#   잔고 확인 후 실제 시장가 매수를 실행합니다.
#
# 실행:
#   UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/dca.sh
#   DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/dca.sh

set -euo pipefail

command -v jq >/dev/null 2>&1 || { echo "jq가 필요합니다." >&2; exit 1; }
command -v bc >/dev/null 2>&1 || { echo "bc가 필요합니다." >&2; exit 1; }

[ -n "${UPBIT_ACCESS_KEY+x}" ] && export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
[ -n "${UPBIT_SECRET_KEY+x}" ] && export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

MARKET="KRW-BTC"
BUY_AMOUNT="5000"   # 1회 매수 금액 (KRW)
TOTAL_ROUNDS=3      # 총 매수 횟수
INTERVAL=5          # 매수 간격 (초)
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
# 1. DCA 설정 확인
# ---------------------------------------------------------------------------

section_setup() {
  local mode_label
  [ "$DRY_RUN" = "false" ] && mode_label="실제 실행" || mode_label="DRY RUN"
  divider "1. DCA 설정 [$mode_label]"

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
      echo "  [오류] $fiat 잔고(${available})가 총 필요 금액(${total_needed})보다 적습니다." >&2
      exit 1
    fi

    printf "  현재가:       %.0f KRW\n" "$current"
    printf "  %s 잔고:  %.0f KRW\n" "$fiat" "$available"
  else
    printf "  현재가:       %.0f KRW\n" "$current"
  fi

  echo "  마켓:         $MARKET"
  printf "  1회 매수금액: %.0f KRW\n" "$BUY_AMOUNT"
  echo "  매수 횟수:    ${TOTAL_ROUNDS}회"
  echo "  매수 간격:    ${INTERVAL}초"
  printf "  총 필요금액:  %.0f KRW\n" "$total_needed"
  [ "$DRY_RUN" != "false" ] && echo "" && echo "  ※ Dry run 모드: 주문이 실행되지 않습니다."
}

# ---------------------------------------------------------------------------
# 2. DCA 매수 실행
# ---------------------------------------------------------------------------

section_execute_dca() {
  divider "2. DCA 매수 실행"

  total_spent=0
  total_volume=0
  success_count=0

  if [ "$DRY_RUN" != "false" ]; then
    echo "  [건너뜀] Dry run 모드이므로 실제 주문을 생략합니다."
    echo ""
    for r in $(seq 1 "$TOTAL_ROUNDS"); do
      current=$(current_price)
      est_volume=$(echo "scale=8; $BUY_AMOUNT / $current" | bc)
      printf "  [%d/%d] 현재가: %.0f KRW  (예상 수량: %s)\n" "$r" "$TOTAL_ROUNDS" "$current" "$est_volume"
      [ "$r" -lt "$TOTAL_ROUNDS" ] && sleep "$INTERVAL"
    done
    return 0
  fi

  for r in $(seq 1 "$TOTAL_ROUNDS"); do
    printf "\n  [%d/%d] 시장가 매수 (%s KRW) ...\n" "$r" "$TOTAL_ROUNDS" "$BUY_AMOUNT"

    order=$(upbit orders create \
      --market "$MARKET" --side "bid" --ord-type "price" --price "$BUY_AMOUNT")
    uuid=$(echo "$order" | jq -r '.uuid')
    echo "  UUID:   $uuid"
    echo "  상태:   $(echo "$order" | jq -r '.state')"

    # 체결 대기
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
      printf "  체결 금액: %.0f KRW\n" "$spent"
      printf "  체결 수량: %.8f\n" "$volume"
      printf "  체결 단가: %.0f KRW\n" "$avg"
      total_spent=$(echo "$total_spent + $spent" | bc)
      total_volume=$(echo "$total_volume + $volume" | bc)
      success_count=$((success_count + 1))
    else
      echo "  체결 상태: $state"
    fi

    [ "$r" -lt "$TOTAL_ROUNDS" ] && echo "  ${INTERVAL}초 후 다음 매수..." && sleep "$INTERVAL"
  done
}

# ---------------------------------------------------------------------------
# 3. 매수 결과 요약
# ---------------------------------------------------------------------------

section_show_summary() {
  divider "3. 매수 결과 요약"

  if [ "$DRY_RUN" != "false" ]; then
    echo "  결과 없음 (Dry run 모드)"
    return 0
  fi

  echo "  성공 횟수:     ${success_count}/${TOTAL_ROUNDS}"
  if (( $(echo "$total_volume > 0" | bc -l) )); then
    avg=$(echo "scale=0; $total_spent / $total_volume" | bc)
    printf "  총 매수금액:   %.0f KRW\n" "$total_spent"
    printf "  총 매수수량:   %.8f\n" "$total_volume"
    printf "  평균 매수단가: %.0f KRW\n" "$avg"
  else
    echo "  성공한 매수가 없습니다."
  fi
}

# ---------------------------------------------------------------------------
# 실행
# ---------------------------------------------------------------------------

section_setup || true
section_execute_dca || true
section_show_summary || true

echo ""
printf '%0.s=' {1..60}; echo ""
echo "  DCA 시나리오 완료"
printf '%0.s=' {1..60}; echo ""
echo ""
echo "  ※ 이 예제는 교육 목적으로 작성되었습니다."
echo "    실제 투자에 사용할 경우 발생하는 손실에 대해 책임지지 않습니다."
