#!/usr/bin/env bash
# 지표 산출 시나리오.
#
# 업비트 Quotation API 응답을 활용하여 투자 지표를 산출하는 예제입니다.
#   - 24시간 누적 거래대금 상위 5개 페어 (KRW 마켓)
#
# Quotation API는 인증이 필요하지 않습니다.
#
# 참고: RSI 산출은 셸 환경의 한계로 포함되지 않습니다.
#       RSI 예제는 python/examples/indicators_kr.py 를 참고하세요.
#
# 실행:
#   bash examples/indicators.sh

set -euo pipefail

command -v jq >/dev/null 2>&1 || { echo "jq가 필요합니다." >&2; exit 1; }

divider() {
  echo ""
  printf '%0.s=' {1..60}; echo ""
  echo "  $1"
  printf '%0.s=' {1..60}; echo ""
}

# ---------------------------------------------------------------------------
# 1. 24시간 누적 거래대금 상위 5개 페어 (KRW 마켓)
# ---------------------------------------------------------------------------

section_top5() {
  divider "1. 24시간 누적 거래대금 상위 5개 (KRW 마켓)"

  tickers=$(upbit tickers list-by-quote-currencies --quote-currencies "KRW")
  total=$(echo "$tickers" | jq 'length')
  echo "  KRW 마켓 수: $total"
  echo ""
  printf "  %-12s %14s %20s\n" "마켓" "현재가" "24h 거래대금"
  echo "  ----------------------------------------------------------------"

  echo "$tickers" | jq -r 'sort_by(-.acc_trade_price_24h | tonumber) | .[:5][] |
    [.market, .trade_price, .acc_trade_price_24h] | @tsv' | \
    awk -F'\t' '{
      acc = $3 + 0
      if (acc >= 1000000000000)
        unit = sprintf("약 %.1f조", acc / 1000000000000)
      else if (acc >= 100000000)
        unit = sprintf("약 %.0f억", acc / 100000000)
      else
        unit = sprintf("%.0f", acc)
      printf "  %-12s %14.0f %20.0f  (%s)\n", $1, $2, acc, unit
    }'
}

# ---------------------------------------------------------------------------
# 실행
# ---------------------------------------------------------------------------

section_top5

echo ""
printf '%0.s=' {1..60}; echo ""
echo "  지표 산출 시나리오 완료"
printf '%0.s=' {1..60}; echo ""
echo ""
echo "  ※ 이 예제는 교육 목적으로 작성되었습니다."
echo "    실제 투자에 사용할 경우 발생하는 손실에 대해 책임지지 않습니다."
