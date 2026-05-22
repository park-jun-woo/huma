# Phase 003: Coverage Gradient Signal

## 목표

hurl PASS/FAIL만으로는 happy path 테스트에서 멈춘다.
핸들러 함수의 미커버 라인을 피드백하여 분기 커버리지를 올린다.

## Phase 001-002와의 차이

| | Phase 001-002 | Phase 003 |
|--|--------------|-----------|
| 판정 | hurl PASS/FAIL | hurl PASS + 커버리지 % |
| 피드백 | "assertion 실패" | "line 41 미커버" |
| 상태 | TODO, PASS | TODO, PASS, IMPROVE, DONE |

## Adapter 인터페이스

```go
type UncoveredLine struct {
    File string
    Line int
    Code string
}

type Adapter interface {
    Build() error
    Start() error
    WaitReady() error
    Stop() error
    Collect(handlerFile string, startLine, endLine int) (covered, total int, uncovered []UncoveredLine, err error)
    Reset() error
}
```

래칫 엔진은 이 인터페이스만 안다. 언어를 모른다.

## Go Adapter

```
Build()  → go build -cover -o ./server ./cmd/server
Start()  → GOCOVERDIR=.hurlfill/coverdata ./server &
Stop()   → SIGTERM → graceful shutdown
Collect()→ go tool covdata textfmt → coverage.out 파싱 → 미커버 라인
Reset()  → coverdata 디렉토리 초기화
```

이 Phase에서는 사용자가 서버를 직접 띄우던 방식 대신,
hurlfill이 계측 서버를 기동/종료한다. hurlfill.yaml에 build/start 설정 추가.

```yaml
server:
  lang: go
  build: "go build -cover -o ./server ./cmd/server"
  start: "./server"
  ready: "http://localhost:8080/health"
  base_url: "http://localhost:8080"
  env:
    DB_DSN: "postgres://localhost:5432/testdb"
```

## IMPROVE / DONE 상태

```
hurl PASS + 커버리지 100%         → PASS
hurl PASS + 커버리지 < 100% (1차) → IMPROVE (미커버 라인 피드백)
hurl PASS + 커버리지 < 100% (2차, 개선 없음) → DONE (best-effort 수용)
```

DONE은 코드의 테스트 가능성 한계를 반영한다.
구체 타입 의존성, 외부 API 호출 등으로 도달 불가능한 분기.

## IMPROVE 출력 예시

```
# IMPROVE  POST /api/users
# Coverage: 65% (11/17)
# UNCOVERED:
#   handler.go:41  if req.Email == ""
#   handler.go:55  if err != nil

## Instructions

1. Read hurl/post_api_users.hurl
2. Add test entries for the uncovered branches above
3. Run `hurlfill next`
```

## 언어 확장

어댑터만 추가하면 다른 언어를 지원한다.

| 언어 | 계측 | 커버리지 수집 |
|------|------|-------------|
| Go | `go build -cover` | `GOCOVERDIR` → `go tool covdata` |
| Python | `coverage run` | `coverage json` |
| Node.js | `NODE_V8_COVERAGE` | `c8 report` |
| Java | JaCoCo agent | `jacoco.exec` |

이 Phase에서는 Go만 구현한다.

## 파일 목록

| 파일 | 상태 |
|------|------|
| `internal/adapter/adapter.go` | 신규 — 인터페이스 |
| `internal/adapter/go.go` | 신규 — Go 어댑터 |
| `internal/coverage/parser.go` | 신규 — coverage.out 파싱 |
| `internal/session/session.go` | 수정 — IMPROVE, DONE 상태 추가 |
| `cmd/next.go` | 수정 — 어댑터 연동 |
| `internal/config/config.go` | 수정 — server.build, server.start 추가 |

## 완료 기준

- [ ] testdata 서버를 계측 빌드하고 hurlfill이 기동/종료한다
- [ ] IMPROVE 출력에 미커버 라인 번호 + 소스 코드가 포함된다
- [ ] 미커버 피드백 후 에이전트가 분기 테스트를 추가하면 커버리지가 올라간다
- [ ] 2회 시도 후 개선 없으면 DONE으로 전환된다
