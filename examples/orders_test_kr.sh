#!/usr/bin/env bash
# 주문 생성 테스트 시나리오.
#
# 주문 생성 테스트 API를 사용하여 실제 주문 없이 다양한 주문 유형을 검증하는 예제입니다.
#   - 주문 가능 정보 확인
#   - 지정가 매수/매도 테스트
#   - 시장가 매수/매도 테스트
#   - 최유리 매수/매도 테스트 (IOC)
#   - 잘못된 주문 요청 검증
#
# 주문 생성 테스트 API는 실제 주문과 동일한 검증 과정을 거치지만,
# 주문이 실제로 생성되지 않으므로 수수료 없이 안전하게 테스트할 수 있습니다.
#
# 실행:
#   UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/orders_test.sh

set -euo pipefail

command -v jq >/dev/null 2>&1 || { echo "jq가 필요합니다." >&2; exit 1; }
command -v bc >/dev/null 2>&1 || { echo "bc가 필요합니다." >&2; exit 1; }

[ -n "${UPBIT_ACCESS_KEY+x}" ] && export UPBIT_ACCESS_KEY="${UPBIT_ACCESS_KEY}"
[ -n "${UPBIT_SECRET_KEY+x}" ] && export UPBIT_SECRET_KEY="${UPBIT_SECRET_KEY}"

MARKET="KRW-BTC"
TEST_VOLUME="0.0001"

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
# 2. 지정가 매수/매도 테스트
# ---------------------------------------------------------------------------

section_test_limit_orders() {
  divider "2. 지정가 매수/매도 테스트"

  tick_size=$(upbit orderbooks list-instruments --markets "$MARKET" | jq -r '.[0].tick_size')
  echo "  호가 단위:     $tick_size"

  best_bid=$(upbit orderbooks list --markets "$MARKET" | jq -r '.[0].orderbook_units[0].bid_price')
  printf "  최고 매수호가: %.0f KRW\n" "$best_bid"

  # 지정가 매수 (3% 할인)
  bid_price=$(echo "scale=0; ($best_bid * 0.97) / $tick_size * $tick_size / 1" | bc)
  echo ""
  echo "  [지정가 매수]"
  printf "  주문 가격:     %.0f KRW (3%% 할인)\n" "$bid_price"
  echo "  주문 수량:     $TEST_VOLUME"
  bid_order=$(upbit orders test-create \
    --market "$MARKET" --side "bid" --ord-type "limit" \
    --price "$bid_price" --volume "$TEST_VOLUME")
  echo "  → UUID:  $(echo "$bid_order" | jq -r '.uuid')"
  echo "  → 상태:  $(echo "$bid_order" | jq -r '.state')"
  echo "  → 유형:  $(echo "$bid_order" | jq -r '.ord_type')"

  # 지정가 매도 (3% 할증)
  ask_price=$(echo "scale=0; ($best_bid * 1.03) / $tick_size * $tick_size / 1" | bc)
  echo ""
  echo "  [지정가 매도]"
  printf "  주문 가격:     %.0f KRW (3%% 할증)\n" "$ask_price"
  echo "  주문 수량:     $TEST_VOLUME"
  ask_order=$(upbit orders test-create \
    --market "$MARKET" --side "ask" --ord-type "limit" \
    --price "$ask_price" --volume "$TEST_VOLUME")
  echo "  → UUID:  $(echo "$ask_order" | jq -r '.uuid')"
  echo "  → 상태:  $(echo "$ask_order" | jq -r '.state')"
  echo "  → 유형:  $(echo "$ask_order" | jq -r '.ord_type')"
}

# ---------------------------------------------------------------------------
# 3. 시장가 매수/매도 테스트
# ---------------------------------------------------------------------------

section_test_market_orders() {
  divider "3. 시장가 매수/매도 테스트"

  min_total=$(upbit orders retrieve-chance --market "$MARKET" | jq -r '.market.bid.min_total')

  # 시장가 매수
  echo "  [시장가 매수]"
  printf "  주문 금액:     %.0f KRW (최소 주문금액)\n" "$min_total"
  buy_order=$(upbit orders test-create \
    --market "$MARKET" --side "bid" --ord-type "price" --price "$min_total")
  echo "  → UUID:  $(echo "$buy_order" | jq -r '.uuid')"
  echo "  → 상태:  $(echo "$buy_order" | jq -r '.state')"
  echo "  → 유형:  $(echo "$buy_order" | jq -r '.ord_type')"

  # 시장가 매도
  echo ""
  echo "  [시장가 매도]"
  echo "  주문 수량:     $TEST_VOLUME"
  sell_order=$(upbit orders test-create \
    --market "$MARKET" --side "ask" --ord-type "market" --volume "$TEST_VOLUME")
  echo "  → UUID:  $(echo "$sell_order" | jq -r '.uuid')"
  echo "  → 상태:  $(echo "$sell_order" | jq -r '.state')"
  echo "  → 유형:  $(echo "$sell_order" | jq -r '.ord_type')"
}

# ---------------------------------------------------------------------------
# 4. 최유리 매수/매도 테스트 (IOC)
# ---------------------------------------------------------------------------

section_test_best_orders() {
  divider "4. 최유리 매수/매도 테스트 (IOC)"

  min_total=$(upbit orders retrieve-chance --market "$MARKET" | jq -r '.market.bid.min_total')

  # 최유리 매수 IOC
  echo "  [최유리 매수 — IOC]"
  printf "  주문 금액:     %.0f KRW (최소 주문금액)\n" "$min_total"
  bid_order=$(upbit orders test-create \
    --market "$MARKET" --side "bid" --ord-type "best" \
    --price "$min_total" --time-in-force "ioc")
  echo "  → UUID:  $(echo "$bid_order" | jq -r '.uuid')"
  echo "  → 상태:  $(echo "$bid_order" | jq -r '.state')"
  echo "  → 유형:  $(echo "$bid_order" | jq -r '.ord_type')"

  # 최유리 매도 IOC
  echo ""
  echo "  [최유리 매도 — IOC]"
  echo "  주문 수량:     $TEST_VOLUME"
  ask_order=$(upbit orders test-create \
    --market "$MARKET" --side "ask" --ord-type "best" \
    --volume "$TEST_VOLUME" --time-in-force "ioc")
  echo "  → UUID:  $(echo "$ask_order" | jq -r '.uuid')"
  echo "  → 상태:  $(echo "$ask_order" | jq -r '.state')"
  echo "  → 유형:  $(echo "$ask_order" | jq -r '.ord_type')"
}

# ---------------------------------------------------------------------------
# 5. 잘못된 주문 요청 검증
# ---------------------------------------------------------------------------

section_test_validation() {
  divider "5. 잘못된 주문 요청 검증"

  # 존재하지 않는 마켓
  echo "  [검증] 존재하지 않는 마켓 (KRW-INVALID)"
  if upbit orders test-create \
    --market "KRW-INVALID" --side "bid" --ord-type "limit" \
    --price "10000" --volume "1" 2>/dev/null; then
    echo "  → 예상과 다르게 성공했습니다."
  else
    echo "  → 에러 발생 (정상)"
  fi

  # 지정가 주문에 가격 누락
  echo ""
  echo "  [검증] 지정가 주문 가격 누락"
  if upbit orders test-create \
    --market "$MARKET" --side "bid" --ord-type "limit" \
    --volume "$TEST_VOLUME" 2>/dev/null; then
    echo "  → 예상과 다르게 성공했습니다."
  else
    echo "  → 에러 발생 (정상)"
  fi
}

# ---------------------------------------------------------------------------
# 실행
# ---------------------------------------------------------------------------

section_check_order_chance
section_test_limit_orders
section_test_market_orders
section_test_best_orders
section_test_validation

echo ""
printf '%0.s=' {1..60}; echo ""
echo "  주문 테스트 시나리오 완료"
printf '%0.s=' {1..60}; echo ""
echo ""
echo "  ※ 주문 생성 테스트 API는 실제 주문을 생성하지 않습니다."
echo "    반환된 UUID는 주문 조회/취소에 사용할 수 없습니다."
