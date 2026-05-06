#!/usr/bin/env bash
# 주문 생성 및 관리 시나리오.
#
# 업비트 Exchange API를 사용하여 주문을 생성하고 관리하는 예제입니다.
#   - 주문 가능 정보 확인
#   - 지정가 매수 주문 생성 → 조회 → 취소
#   - 보유 자산 확인 + 시장가 매도 가능 여부
#   - 완료 주문 목록 조회
#
# Dry run (기본):
#   주문 관련 정보 조회만 수행하며, 실제 주문을 생성하지 않습니다.
#
# 실제 실행 (DRY_RUN=false):
#   지정가 매수 주문 생성 → 조회 → 취소를 실행합니다.
#
# 실행:
#   UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/orders.sh
#   DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/orders.sh

set -euo pipefail

command -v jq >/dev/null 2>&1 || { echo "jq가 필요합니다." >&2; exit 1; }
command -v bc >/dev/null 2>&1 || { echo "bc가 필요합니다." >&2; exit 1; }

[ -n "${UPBIT_ACCESS_KEY+x}" ] && export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
[ -n "${UPBIT_SECRET_KEY+x}" ] && export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

MARKET="KRW-BTC"
DRY_RUN="${DRY_RUN:-true}"

divider() {
  echo ""
  printf '%0.s=' {1..60}; echo ""
  echo "  $1"
  printf '%0.s=' {1..60}; echo ""
}

# ---------------------------------------------------------------------------
# 1. 주문 가능 정보 확인
# ---------------------------------------------------------------------------

section_check_order_chance() {
  divider "1. 주문 가능 정보 확인"

  chance=$(upbit orders retrieve-chance --market "$MARKET")
  printf "  마켓:          %s\n"   "$MARKET"
  printf "  매수 수수료:   %s\n"   "$(echo "$chance" | jq -r '.bid_fee')"
  printf "  매도 수수료:   %s\n"   "$(echo "$chance" | jq -r '.ask_fee')"
  printf "  매수 주문유형: %s\n"   "$(echo "$chance" | jq -r '.market.bid_types | @csv')"
  printf "  매도 주문유형: %s\n"   "$(echo "$chance" | jq -r '.market.ask_types | @csv')"
  printf "  최소 매수금액: %s\n"   "$(echo "$chance" | jq -r '.market.bid.min_total')"
  printf "  최소 매도금액: %s\n"   "$(echo "$chance" | jq -r '.market.ask.min_total')"
}

# ---------------------------------------------------------------------------
# 2. 지정가 매수 → 조회 → 취소
# ---------------------------------------------------------------------------

section_limit_bid_and_cancel() {
  divider "2. 지정가 매수 → 조회 → 취소"

  # 호가 단위 조회
  tick_size=$(upbit orderbooks list-instruments --markets "$MARKET" | jq -r '.[0].tick_size')
  echo "  호가 단위:     $tick_size"

  # 최고 매수호가 조회
  best_bid=$(upbit orderbooks list --markets "$MARKET" | jq -r '.[0].orderbook_units[0].bid_price')
  printf "  최고 매수호가: %.0f KRW\n" "$best_bid"

  # 3% 낮은 가격 계산 (호가 단위 적용)
  target_price=$(echo "scale=0; ($best_bid * 0.97) / $tick_size * $tick_size / 1" | bc)
  volume="0.0001"
  printf "  주문 가격:     %'.0f KRW (3%% 할인)\n" "$target_price"
  echo "  주문 수량:     $volume"

  if [ "$DRY_RUN" != "false" ]; then
    echo ""
    echo "  [건너뜀] Dry run 모드이므로 실제 주문을 생략합니다."
    return 0
  fi

  # 주문 생성
  order=$(upbit orders create \
    --market "$MARKET" \
    --side "bid" \
    --ord-type "limit" \
    --price "$target_price" \
    --volume "$volume")
  uuid=$(echo "$order" | jq -r '.uuid')
  echo ""
  echo "  주문 생성 완료"
  echo "  UUID:   $uuid"
  echo "  상태:   $(echo "$order" | jq -r '.state')"

  # 주문 조회
  info=$(upbit orders retrieve --uuid "$uuid")
  echo "  조회:   상태=$(echo "$info" | jq -r '.state'), 체결수량=$(echo "$info" | jq -r '.executed_volume')"

  # 주문 취소
  upbit orders cancel --uuid "$uuid" > /dev/null
  for _ in $(seq 1 10); do
    cancelled=$(upbit orders retrieve --uuid "$uuid")
    state=$(echo "$cancelled" | jq -r '.state')
    if [ "$state" = "done" ] || [ "$state" = "cancel" ]; then
      break
    fi
    sleep 0.5
  done
  echo "  취소:   상태=$state"
}

# ---------------------------------------------------------------------------
# 3. 보유 자산 확인 + 시장가 매도 가능 여부
# ---------------------------------------------------------------------------

section_check_market_sell() {
  divider "3. 보유 자산 확인 + 시장가 매도 가능 여부"

  coin=$(echo "$MARKET" | cut -d'-' -f2)

  accounts=$(upbit accounts list)
  target=$(echo "$accounts" | jq --arg c "$coin" '.[] | select(.currency == $c)')

  if [ -z "$target" ] || [ "$(echo "$target" | jq -r '.balance')" = "0" ]; then
    echo "  $coin 보유량이 없습니다."
  else
    balance=$(echo "$target" | jq -r '.balance')
    avg_buy=$(echo "$target" | jq -r '.avg_buy_price')
    echo "  보유 자산:   $coin"
    printf "  보유 수량:   %.8f\n" "$balance"
    printf "  평균 매수가: %.0f KRW\n" "$avg_buy"
  fi

  chance=$(upbit orders retrieve-chance --market "$MARKET")
  ask_types=$(echo "$chance" | jq -r '.market.ask_types | @csv')
  if echo "$ask_types" | grep -q "market"; then
    echo "  시장가 매도: 가능"
  else
    echo "  시장가 매도: 불가능 (ask_types=$ask_types)"
  fi
}

# ---------------------------------------------------------------------------
# 4. 완료 주문 목록 조회
# ---------------------------------------------------------------------------

section_list_closed_orders() {
  divider "4. 완료 주문 목록 조회"

  closed=$(upbit orders list-closed --market "$MARKET" --limit 5)
  count=$(echo "$closed" | jq 'length')

  if [ "$count" = "0" ]; then
    echo "  완료 주문이 없습니다."
    return 0
  fi

  printf "  %-10s %4s %8s %6s %14s\n" "UUID" "구분" "유형" "상태" "체결수량"
  echo "  --------------------------------------------------"
  echo "$closed" | jq -r '.[] | [.uuid[:8], .side, .ord_type, .state, .executed_volume] | @tsv' | \
    awk -F'\t' '{
      side = ($2 == "bid") ? "매수" : "매도"
      printf "  %s... %4s %8s %6s %14s\n", $1, side, $3, $4, $5
    }'
}

# ---------------------------------------------------------------------------
# 실행
# ---------------------------------------------------------------------------

section_check_order_chance || true
section_limit_bid_and_cancel || true
section_check_market_sell || true
section_list_closed_orders || true

echo ""
printf '%0.s=' {1..60}; echo ""
echo "  주문 시나리오 완료"
printf '%0.s=' {1..60}; echo ""
echo ""
echo "  ※ 이 예제는 교육 목적으로 작성되었습니다."
echo "    실제 투자에 사용할 경우 발생하는 손실에 대해 책임지지 않습니다."
