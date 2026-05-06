[English](./README.md) | 한국어

# Upbit CLI

[Upbit REST API](https://docs.upbit.com)를 위한 공식 CLI입니다.

<!-- x-release-please-start-version -->

## 설치

### npm으로 설치하기

```sh
npm install -g @upbit-official/upbit-cli
```

### Go로 설치하기

CLI를 로컬에서 테스트하거나 설치하려면 [Go](https://go.dev/doc/install) 1.22 이상이 필요합니다.

```sh
go install 'github.com/upbit-official/upbit-cli/cmd/upbit@latest'
```

`go install` 실행 후 바이너리는 Go bin 디렉터리에 설치됩니다.

- **기본 위치**: `$HOME/go/bin` (또는 `GOPATH`가 설정된 경우 `$GOPATH/bin`)
- **경로 확인**: `go env GOPATH`로 base 디렉터리를 확인할 수 있습니다.

설치 후에도 명령이 실행되지 않는다면 Go bin 디렉터리를 PATH에 추가하세요.

```sh
# 셸 프로파일(.zshrc, .bashrc 등)에 추가
export PATH="$PATH:$(go env GOPATH)/bin"
```

<!-- x-release-please-end -->

### 로컬에서 실행

이 프로젝트를 clone한 뒤 `scripts/run` 스크립트로 CLI를 로컬에서 실행할 수 있습니다.

```sh
./scripts/run args...
```

## 사용법

CLI는 리소스 기반 명령 구조를 따릅니다.

```sh
upbit [resource] <command> [flags...]
```

```sh
upbit accounts list \
  --access-key "$UPBIT_ACCESS_KEY" \
  --secret-key "$UPBIT_SECRET_KEY"
```

각 명령의 자세한 옵션은 `--help` 플래그로 확인하세요.

실행 가능한 예제는 [`examples/`](examples/) 디렉터리에서 확인할 수 있습니다.

### 환경 변수

| 환경 변수 | 설명 | 필수 여부 | 기본값 |
| --- | --- | --- | --- |
| `UPBIT_ACCESS_KEY` | Upbit API 인증에 사용하는 액세스 키. 자세한 내용은 https://docs.upbit.com/reference/auth 를 참고하세요. | 아니오 | `null` |
| `UPBIT_SECRET_KEY` | API 요청 서명에 사용하는 시크릿 키. 자세한 내용은 https://docs.upbit.com/reference/auth 를 참고하세요. | 아니오 | `null` |

### 전역 플래그

- `--access-key` — Upbit API 인증에 사용하는 액세스 키. 자세한 내용은 https://docs.upbit.com/reference/auth 를 참고하세요.
  (`UPBIT_ACCESS_KEY` 환경 변수로도 설정 가능)
- `--secret-key` — API 요청 서명에 사용하는 시크릿 키. 자세한 내용은 https://docs.upbit.com/reference/auth 를 참고하세요.
  (`UPBIT_SECRET_KEY` 환경 변수로도 설정 가능)
- `--help` — 명령어 사용법 표시
- `--debug` — 디버그 로깅 활성화 (HTTP 요청/응답 세부 정보 포함)
- `--version`, `-v` — CLI 버전 표시
- `--base-url` — 커스텀 API 백엔드 URL 사용
- `--format` — 출력 형식 변경 (`auto`, `explore`, `json`, `jsonl`, `pretty`, `raw`, `yaml`)
- `--format-error` — 오류 출력 형식 변경 (`auto`, `explore`, `json`, `jsonl`, `pretty`, `raw`, `yaml`)
- `--transform` — [GJSON 구문](https://github.com/tidwall/gjson/blob/master/SYNTAX.md)으로 데이터 출력 변환
- `--transform-error` — [GJSON 구문](https://github.com/tidwall/gjson/blob/master/SYNTAX.md)으로 오류 출력 변환

## Go SDK 버전 연결

개발 목적으로 CLI를 다른 버전의 Upbit Go SDK에 연결하려면 `./scripts/link` 스크립트를 사용할 수 있습니다.

저장소의 특정 버전(브랜치, git 태그, 커밋 해시)에 연결하려면:

```bash
./scripts/link github.com/org/repo@version
```

로컬 SDK 사본에 연결하려면:

```bash
./scripts/link ../path/to/upbit-go
```

인자 없이 실행하면 기본값은 `../upbit-go`입니다.

## 버전 관리

이 패키지는 일반적으로 [SemVer](https://semver.org/spec/v2.0.0.html) 규약을 따르지만, 일부 하위 호환성에 영향을 줄 수 있는 변경이 마이너 버전에 포함될 수 있습니다.

1. 기술적으로는 공개되어 있지만 외부 사용을 의도하지 않았거나 문서화되지 않은 CLI 내부 변경
2. 대부분의 사용자에게 실질적인 영향이 없을 것으로 판단되는 변경

하위 호환성은 중요하게 고려하고 있으며, 원활한 업그레이드 경험을 제공하기 위해 노력하고 있습니다.

피드백을 환영합니다. 질문, 버그 제보, 개선 제안이 있다면 open-api@upbit.com 으로 연락해 주세요.

## 기여

현재 업비트 CLI는 초기 출시 단계로, 외부 Issue 등록 및 PR 기여는 아직 운영하고 있지 않습니다.
버그 제보나 개선 의견은 개발자 지원 채널(open-api@upbit.com)을 통해 전달해 주시기 바랍니다.
외부 기여 채널은 향후 운영 안정화 상황에 맞춰 순차적으로 오픈하는 방안을 검토하고 있습니다.

© 2026 Dunamu Inc. All rights reserved.
