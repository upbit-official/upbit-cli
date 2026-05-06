#!/usr/bin/env bash
# TP/SL (익절/손절) 자동 매도 시나리오.
#
# 업비트 Exchange API를 사용하여 현재가를 폴링하고,
# 목표가(익절) 또는 손절가에 도달하면 시장가 매도를 실행하는 예제입니다.
#   - 현재가 조회 및 TP/SL 가격 계산
#   - 가격 모니터링 (REST 폴링)
#   - 시장가 매도 실행
#
# 업비트 API는 예약주문(스톱/리밋)을 지원하지 않으므로
# 클라이언트 사이드에서 폴링 방식으로 구현합니다.
#
# Dry run (기본):
#   현재가 기준, 주문 실행 안 함.
#
# 실제 실행 (DRY_RUN=false):
#   평균 매수가 기준, 실제 매도 주문 실행. 반드시 보유 자산이 있어야 합니다.
#
# 실행:
#   UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/tp_sl.sh
#   DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/tp_sl.sh

set -euo pipefail

command -v jq >/dev/null 2>&1 || { echo "jq가 필요합니다." >&2; exit 1; }
command -v bc >/dev/null 2>&1 || { echo "bc가 필요합니다." >&2; exit 1; }

[ -n "${UPBIT_ACCESS_KEY+x}" ] && export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
[ -n "${UPBIT_SECRET_KEY+x}" ] && export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

MARKET="KRW-BTC"
SELL_VOLUME="0.0001"
TP_PERCENT="3"    # 익절 기준: +3%
SL_PERCENT="2"    # 손절 기준: -2%
POLL_INTERVAL=1   # 폴링 간격 (초)
MAX_POLLS=10      # 최대 폴링 횟수
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
# 1. TP/SL 설정
# ---------------------------------------------------------------------------

section_setup() {
  local mode_label
  [ "$DRY_RUN" = "false" ] && mode_label="실제 실행" || mode_label="DRY RUN"
  divider "1. TP/SL 설정 [$mode_label]"

  local current
  current=$(current_price)

  if [ "$DRY_RUN" = "false" ]; then
    coin=$(echo "$MARKET" | cut -d'-' -f2)
    accounts=$(upbit accounts list)
    target=$(echo "$accounts" | jq --arg c "$coin" '.[] | select(.currency == $c)')

    if [ -z "$target" ]; then
      echo "  [오류] $coin 보유량이 없습니다." >&2; exit 1
    fi

    balance=$(echo "$target" | jq -r '.balance')
    if (( $(echo "$balance < $SELL_VOLUME" | bc -l) )); then
      echo "  [오류] $coin 보유량(${balance})이 매도 수량(${SELL_VOLUME})보다 적습니다." >&2; exit 1
    fi

    BASE_PRICE=$(echo "$target" | jq -r '.avg_buy_price')
    echo "  잔고:       $coin $balance"
    printf "  기준가:     %.0f KRW (평균 매수가)\n" "$BASE_PRICE"
    printf "  현재가:     %.0f KRW\n" "$current"
  else
    BASE_PRICE="$current"
    printf "  기준가:     %.0f KRW (현재가)\n" "$BASE_PRICE"
  fi

  TP_PRICE=$(echo "scale=0; $BASE_PRICE * (1 + $TP_PERCENT / 100) / 1" | bc)
  SL_PRICE=$(echo "scale=0; $BASE_PRICE * (1 - $SL_PERCENT / 100) / 1" | bc)

  echo "  마켓:       $MARKET"
  echo "  매도 수량:  $SELL_VOLUME"
  printf "  익절 가격:  %.0f KRW (+%s%%)\n" "$TP_PRICE" "$TP_PERCENT"
  printf "  손절 가격:  %.0f KRW (-%s%%)\n" "$SL_PRICE" "$SL_PERCENT"
  echo "  폴링 간격:  ${POLL_INTERVAL}초 × 최대 ${MAX_POLLS}회"
  [ "$DRY_RUN" != "false" ] && echo "" && echo "  ※ Dry run 모드: 주문이 실행되지 않습니다."
}

# ---------------------------------------------------------------------------
# 2. 가격 모니터링
# ---------------------------------------------------------------------------

section_monitor_price() {
  divider "2. 가격 모니터링 시작"

  TRIGGER="timeout"

  for poll in $(seq 1 "$MAX_POLLS"); do
    current=$(current_price)
    printf "  [%3d/%d] 현재가: %.0f KRW" "$poll" "$MAX_POLLS" "$current"

    if (( $(echo "$current >= $TP_PRICE" | bc -l) )); then
      echo "  → 익절 도달!"
      TRIGGER="tp"; break
    fi
    if (( $(echo "$current <= $SL_PRICE" | bc -l) )); then
      echo "  → 손절 도달!"
      TRIGGER="sl"; break
    fi

    echo ""
    [ "$poll" -lt "$MAX_POLLS" ] && sleep "$POLL_INTERVAL"
  done

  if [ "$TRIGGER" = "timeout" ]; then
    echo ""
    echo "  최대 폴링 횟수(${MAX_POLLS}회)에 도달했습니다."
  fi
}

# ---------------------------------------------------------------------------
# 3. 시장가 매도 실행
# ---------------------------------------------------------------------------

section_execute_sell() {
  local label
  [ "$TRIGGER" = "tp" ] && label="익절" || label="손절"
  divider "3. 매도 실행 (트리거: $label)"

  if [ "$DRY_RUN" != "false" ]; then
    echo "  [건너뜀] Dry run 모드이므로 실제 주문을 생략합니다."
    return 0
  fi

  order=$(upbit orders create \
    --market "$MARKET" --side "ask" --ord-type "market" --volume "$SELL_VOLUME")
  uuid=$(echo "$order" | jq -r '.uuid')
  echo "  주문 생성 완료"
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

  echo "  체결 상태: $state"
  echo "  체결 수량: $(echo "$info" | jq -r '.executed_volume')"
  if [ "$state" = "done" ]; then
    total=$(echo "$info" | jq '[.trades[].funds | tonumber] | add // 0')
    printf "  체결 금액: %.0f KRW\n" "$total"
  fi
}

# ---------------------------------------------------------------------------
# 실행
# ---------------------------------------------------------------------------

section_setup || true
section_monitor_price || true

if [ "$TRIGGER" = "timeout" ]; then
  echo ""
  echo "  시간 초과로 매도를 실행하지 않습니다."
else
  section_execute_sell
fi

echo ""
printf '%0.s=' {1..60}; echo ""
echo "  TP/SL 시나리오 완료"
printf '%0.s=' {1..60}; echo ""
echo ""
echo "  ※ 이 예제는 교육 목적으로 작성되었습니다."
echo "    실제 투자에 사용할 경우 발생하는 손실에 대해 책임지지 않습니다."
