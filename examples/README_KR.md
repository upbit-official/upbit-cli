[English](./README.md) | 한국어

# 예제 코드

이 디렉토리에는 업비트 CLI SDK의 주요 기능을 시나리오 기반으로 설명하는 예제가 포함되어 있습니다.

각 예제는 실제 거래 워크플로우를 기반으로 작성되었으며, **Dry run 모드**가 기본값입니다.
Dry run 모드에서는 조회 작업만 수행하며, 주문 생성·출금 등 쓰기 작업은 생략됩니다.

> **참고:** 엔드포인트별 자동 생성 예제는 기능별 하위 디렉토리에 정리되어 있습니다(예: [`examples/orders/`](./orders/)). 영문 시나리오 예제는 [`README.md`](README.md)를 참고하세요.

---

## 빠른 시작

처음이라면 아래 순서로 예제를 살펴보는 것을 권장합니다:

1. [`quotation_kr.sh`](./quotation_kr.sh) — 인증 없이 시세·체결·호가 데이터를 빠르게 조회
2. [`indicators_kr.sh`](./indicators_kr.sh) — Quotation 데이터로 거래대금 상위 페어 확인
3. [`orders_kr.sh`](./orders_kr.sh) — Dry run으로 주문 생성·조회·취소 흐름을 먼저 확인

### 인증이 필요 없는 예제

아래 예제는 API 키 없이 즉시 실행 가능합니다:
- [`quotation_kr.sh`](./quotation_kr.sh)
- [`indicators_kr.sh`](./indicators_kr.sh)

### 인증이 필요한 예제

아래 예제는 업비트 API 키가 필요합니다:
- [`orders_kr.sh`](./orders_kr.sh)
- [`orders_test_kr.sh`](./orders_test_kr.sh)
- [`deposits_kr.sh`](./deposits_kr.sh)
- [`withdrawals_kr.sh`](./withdrawals_kr.sh)
- [`dca_kr.sh`](./dca_kr.sh)
- [`tp_sl_kr.sh`](./tp_sl_kr.sh)

### 사전 준비

```bash
# jq 설치 (JSON 파싱에 필요)
# macOS
brew install jq
# Ubuntu/Debian
apt-get install jq

# 업비트 API 키 발급 후 환경변수 설정
export UPBIT_ACCESS_KEY=<your-access-key>
export UPBIT_SECRET_KEY=<your-secret-key>
```

### 실행 전 참고사항

- **조회 동작을 먼저 확인하고 싶다면** `quotation_kr.sh`, `indicators_kr.sh`로 시작하세요.
- **주문·출금이 포함된 예제를 검토할 때는** Dry run 모드에서 동작을 먼저 확인하세요.
- **`DRY_RUN=false`로 실행하면** 주문 생성·출금·자동 매매가 실제로 실행됩니다.

---

## 예제 목록

| 예제 | 목적 | 인증 | 기본 동작 |
|---|---|---|---|
| [`quotation_kr.sh`](./quotation_kr.sh) | 시세·캔들·체결·호가 조회 | 불필요 | 안전한 조회 전용 |
| [`indicators_kr.sh`](./indicators_kr.sh) | 24h 거래대금 상위 페어 조회 | 불필요 | 안전한 조회 전용 |
| [`orders_kr.sh`](./orders_kr.sh) | 주문 생성·조회·취소 흐름 | 필요 | Dry run 기본; 실제 주문 가능 |
| [`orders_test_kr.sh`](./orders_test_kr.sh) | 주문 생성 테스트 API로 모든 주문 유형 검증 | 필요 | 항상 안전 — 실제 주문 없음 |
| [`deposits_kr.sh`](./deposits_kr.sh) | 입금 주소 및 내역 관리 | 필요 | Dry run 기본; 일부 작업 실제 효과 가능 |
| [`withdrawals_kr.sh`](./withdrawals_kr.sh) | 출금 정보 확인 및 출금 흐름 | 필요 | Dry run 기본; 실제 출금 주의 |
| [`dca_kr.sh`](./dca_kr.sh) | 정기 시장가 매수(DCA) 자동화 | 필요 | Dry run 기본; 실제 주문 가능 |
| [`tp_sl_kr.sh`](./tp_sl_kr.sh) | 익절/손절 자동 매도 | 필요 | Dry run 기본; 실제 주문 가능 |

---

## 시나리오 예제 상세

### `quotation_kr.sh` — 시세 조회 입문

인증 없이 사용 가능한 Quotation API를 빠르게 체험할 수 있습니다.

**주요 기능:**
- KRW 마켓 목록 조회 및 유의 종목 필터링
- 현재가(티커) 조회 — 단일, 복수, 마켓 단위
- 캔들(OHLCV) 데이터 조회 — 5분봉·일봉·주봉·월봉
- 최근 체결 내역 조회 — `days_ago` 파라미터 활용
- 호가 데이터 조회

```bash
bash examples/quotation_kr.sh
```

---

### `orders_kr.sh` — 주문 생성 및 관리

지정가 매수 주문 생성부터 조회·취소까지 전체 흐름을 다룹니다.

```bash
# Dry run (기본) — 조회만, 주문 없음
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/orders_kr.sh

# 실제 실행 — 지정가 매수 → 조회 → 취소
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/orders_kr.sh
```

---

### `orders_test_kr.sh` — 주문 유형 검증 (test-create API)

주문 생성 테스트 API를 사용하여 모든 주문 유형을 안전하게 검증합니다.
실제 주문이 생성되지 않으므로 수수료 없이 테스트할 수 있습니다.

```bash
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/orders_test_kr.sh
```

---

### `deposits_kr.sh` — 입금 주소 및 내역 관리

```bash
# Dry run (기본) — 조회만
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/deposits_kr.sh

# 실제 실행 — 입금 주소 생성 포함
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/deposits_kr.sh
```

---

### `withdrawals_kr.sh` — 출금 관리

> **경고:** 출금은 비가역적입니다. 항상 주소와 금액을 확인한 뒤 실행하세요.

```bash
# Dry run (기본) — 조회만
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/withdrawals_kr.sh
```

---

### `dca_kr.sh` — DCA 정기 매수 자동화

고정 금액을 반복 시장가 매수하는 DCA 전략을 구현합니다.

```bash
# Dry run (기본) — 현재가만 조회, 주문 없음
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/dca_kr.sh

# 실제 실행
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/dca_kr.sh
```

---

### `tp_sl_kr.sh` — 익절/손절 자동 매도

현재가를 폴링하여 목표가 또는 손절가에 도달하면 시장가 매도를 실행합니다.

```bash
# Dry run (기본) — 모니터링만, 주문 없음
UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/tp_sl_kr.sh

# 실제 실행 — 평균 매수가 기준, 자동 매도
DRY_RUN=false UPBIT_ACCESS_KEY=<key> UPBIT_SECRET_KEY=<secret> bash examples/tp_sl_kr.sh
```

---

### `indicators_kr.sh` — 투자 지표

KRW 마켓 24시간 거래대금 상위 페어를 조회합니다.

```bash
bash examples/indicators_kr.sh
```

> RSI 산출은 셸 환경의 한계로 포함되지 않습니다. RSI 예제는 [`python/examples/indicators_kr.py`](https://github.com/upbit-official/upbit-sdk-python/blob/main/examples/indicators_kr.py)를 참고하세요.
