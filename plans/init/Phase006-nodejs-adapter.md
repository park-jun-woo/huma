# Phase 006: Node.js Coverage Adapter

## 목표

Node.js 백엔드(Express, Fastify, NestJS, Koa)에 대한 커버리지 어댑터를 구현한다.

## 라우트 추출 (사용자 몫)

| 프레임워크 | 방법 |
|-----------|------|
| Express | `express-list-endpoints` 패키지 |
| Fastify | `fastify.printRoutes()` |
| NestJS | `@nestjs/swagger` → OpenAPI |
| Koa | `koa-router` inspect |

## 커버리지 어댑터

Node.js는 V8 내장 커버리지를 사용한다. 프레임워크 무관.

### hurlfill.yaml

```yaml
base_url: "http://localhost:3000"
hurl_dir: "hurl"
server:
  lang: node
  build: "npm install"
  start: "node server.js"
  ready: "http://localhost:3000/health"
```

### NodeAdapter

```go
type NodeAdapter struct {
    cfg      *config.ServerConfig
    baseURL  string
    coverDir string   // .hurlfill/v8cov
    proc     *exec.Cmd
}
```

| 메서드 | 동작 |
|--------|------|
| `Build()` | `npm install` (한 번) |
| `Start()` | `NODE_V8_COVERAGE=.hurlfill/v8cov node server.js` |
| `WaitReady()` | ready URL 폴링 (GoAdapter와 동일 로직 재사용) |
| `Stop()` | SIGINT → V8이 커버리지 JSON 덤프 |
| `Collect()` | `npx c8 report --reporter=json --temp-directory=.hurlfill/v8cov` → 파싱 |
| `Reset()` | v8cov 디렉토리 초기화 |

### Stop() 동작 원리

Node.js는 `NODE_V8_COVERAGE` 환경변수가 설정되면, 프로세스 종료 시
해당 디렉토리에 V8 커버리지 JSON을 자동 덤프한다.
SIGINT → graceful shutdown → V8 coverage flush → 파일 생성.

### c8 report JSON 포맷 (istanbul 호환)

`c8 report --reporter=json`은 raw V8 포맷이 아니라 istanbul 호환 JSON을 출력한다.
라인 기반이므로 바이트 오프셋 → 라인 변환이 불필요하다.

```json
{
  "/app/routes/user.js": {
    "path": "/app/routes/user.js",
    "statementMap": {
      "0": { "start": { "line": 5, "column": 0 }, "end": { "line": 5, "column": 30 } },
      "1": { "start": { "line": 10, "column": 4 }, "end": { "line": 10, "column": 25 } }
    },
    "s": {
      "0": 1,
      "1": 0
    }
  }
}
```

`s` 맵에서 count가 0인 statement의 `statementMap` 라인을 추출하면 미커버 라인.

## 파일 목록

| 파일 | 상태 |
|------|------|
| `internal/adapter/node_adapter.go` | 신규 — NodeAdapter 타입 |
| `internal/adapter/new_node_adapter.go` | 신규 — 생성자 |
| `internal/adapter/node_build.go` | 신규 — npm install |
| `internal/adapter/node_start.go` | 신규 — NODE_V8_COVERAGE 설정 + 서버 기동 |
| `internal/adapter/node_stop.go` | 신규 — SIGINT |
| `internal/adapter/node_collect.go` | 신규 — c8 report + istanbul JSON 파싱 |
| `internal/adapter/node_reset.go` | 신규 — v8cov 디렉토리 초기화 |
| `internal/coverage/parse_istanbul.go` | 신규 — istanbul JSON 파서 |

## 완료 기준

- [ ] NodeAdapter가 Adapter 인터페이스를 구현한다
- [ ] istanbul JSON에서 미커버 라인이 추출된다
- [ ] lang=node일 때 next/verify가 NodeAdapter를 사용한다
- [ ] go test 전부 통과
- [ ] filefunc validate 0 violations
