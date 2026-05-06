#!/usr/bin/env bash
# 시세 조회 시나리오.
#
# 업비트 Quotation API를 사용하여 시세 데이터를 조회하는 예제입니다.
#   - 마켓 목록 조회 (KRW 페어)
#   - 현재가(티커) 조회 — 단일/복수/마켓 단위
#   - 캔들(OHLCV) 조회 — 분/일/주/월, to 파라미터
#   - 최근 체결 내역 조회 (days_ago 포함)
#   - 호가 조회
#
# Quotation API는 인증이 필요하지 않습니다.
#
# 실행:
#     bash examples/quotation_kr.sh

set -euo pipefail

# ---------------------------------------------------------------------------
# 헬퍼
# ---------------------------------------------------------------------------
divider() {
  echo ""
  echo "============================================================"
  echo "  $1"
  echo "============================================================"
}

# ---------------------------------------------------------------------------
# 1. 마켓 목록 조회
# ---------------------------------------------------------------------------
list_markets() {
  divider "1. 마켓 목록"

  echo "  [전체 마켓]"
  upbit trading-pairs list --is-details=false | jq -r '
    "  전체 마켓 수: \(length)\n",
    "  \("마켓" | . + " " * (12 - length))\("한글명" | . + " " * (12 - length))영문명",
    "  --------------------------------------------------",
    (.[:5][] | "  \(.market | . + " " * (12 - length))\(.korean_name | . + " " * (12 - length))\(.english_name)")
  '

  echo ""
  echo "  [유의 종목]"
  upbit trading-pairs list --is-details=true | jq -r '
    [.[] | select(.market_warning == "CAUTION")] |
    "  유의 종목 수: \(length)"
  '
}

# ---------------------------------------------------------------------------
# 2. 현재가(티커) 조회
# ---------------------------------------------------------------------------
get_ticker() {
  divider "2. 현재가(티커) 조회"

  echo "  [단일 마켓]"
  upbit tickers list-by-trading-pairs --markets "KRW-BTC" | jq -r '
    .[0] |
    "  마켓:          \(.market)",
    "  현재가:        \(.trade_price) KRW",
    "  24h 거래대금:  \(.acc_trade_price_24h) KRW",
    "  변동률:        \(.signed_change_rate)"
  '

  echo ""
  echo "  [복수 마켓]"
  upbit tickers list-by-trading-pairs --markets "KRW-BTC,KRW-ETH" | jq -r '
    .[] | "  \(.market | . + " " * (12 - length))현재가: \(.trade_price)"
  '

  echo ""
  echo "  [전체 KRW 티커]"
  upbit tickers list-by-quote-currencies --quote-currencies "KRW" | jq -r '
    "  KRW 마켓 티커 수: \(length)"
  '
}

# ---------------------------------------------------------------------------
# 3. 캔들 조회
# ---------------------------------------------------------------------------
get_candles() {
  divider "3. 캔들 조회"

  echo "  [5분봉]"
  upbit candles list-minutes --unit 5 --market "KRW-BTC" --count 3 | jq -r '
    if length == 0 then "  캔들 데이터가 없습니다."
    else
      (.[] | "  \(.candle_date_time_kst)  시가: \(.opening_price)  종가: \(.trade_price)  거래량: \(.candle_acc_trade_volume)")
    end
  '

  echo ""
  echo "  [to 파라미터 — 2025-01-01 이전 1분봉]"
  upbit candles list-minutes --unit 1 --market "KRW-BTC" --to "2025-01-01T00:00:00Z" --count 1 | jq -r '
    if length == 0 then "  지정한 기간의 데이터가 없습니다."
    else .[0] | "  \(.candle_date_time_kst)  종가: \(.trade_price)"
    end
  '

  echo ""
  echo "  [일봉]"
  upbit candles list-days --market "KRW-BTC" --count 3 | jq -r '
    if length == 0 then "  일봉 데이터가 없습니다."
    else .[] | "  \(.candle_date_time_kst[:10])  종가: \(.trade_price)"
    end
  '

  echo ""
  echo "  [주봉]"
  upbit candles list-weeks --market "KRW-BTC" --count 3 | jq -r '
    if length == 0 then "  주봉 데이터가 없습니다."
    else .[] | "  \(.candle_date_time_kst[:10])  종가: \(.trade_price)"
    end
  '

  echo ""
  echo "  [월봉]"
  upbit candles list-months --market "KRW-BTC" --count 3 | jq -r '
    if length == 0 then "  월봉 데이터가 없습니다."
    else .[] | "  \(.candle_date_time_kst[:10])  종가: \(.trade_price)"
    end
  '
}

# ---------------------------------------------------------------------------
# 4. 최근 체결 조회
# ---------------------------------------------------------------------------
get_trades() {
  divider "4. 최근 체결"

  upbit trades list --market "KRW-BTC" --count 5 | jq -r '
    .[] |
    (if .ask_bid == "BID" then "매수" else "매도" end) as $side |
    "  \(.trade_date_utc) \(.trade_time_utc)  \($side)  \(.trade_price)  \(.trade_volume)"
  '

  echo ""
  echo "  [1일 전 체결]"
  upbit trades list --market "KRW-BTC" --days-ago 1 --count 3 | jq -r '
    if length == 0 then "  1일 전 체결 내역이 없습니다."
    else .[] |
      (if .ask_bid == "BID" then "매수" else "매도" end) as $side |
      "  \(.trade_date_utc) \(.trade_time_utc)  \($side)  \(.trade_price)"
    end
  '
}

# ---------------------------------------------------------------------------
# 5. 호가 조회
# ---------------------------------------------------------------------------
get_orderbook() {
  divider "5. 호가"

  upbit orderbooks list --markets "KRW-BTC" | jq -r '
    .[0] |
    "  총 매도 잔량: \(.total_ask_size)",
    "  총 매수 잔량: \(.total_bid_size)",
    "",
    (.orderbook_units[:5][] |
      "  매도: \(.ask_price) (\(.ask_size))   매수: \(.bid_price) (\(.bid_size))"
    )
  '
}

# ---------------------------------------------------------------------------
# 실행
# ---------------------------------------------------------------------------
list_markets
get_ticker
get_candles
get_trades
get_orderbook

echo ""
echo "============================================================"
echo "  시세 조회 시나리오 완료"
echo "============================================================"
echo ""
echo "  ※ 이 예제는 교육 목적으로 작성되었습니다."
echo "    실제 투자에 사용할 경우 발생하는 손실에 대해 책임지지 않습니다."
