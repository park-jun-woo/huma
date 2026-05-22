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
| `WaitReady()` | ready URL 폴링 |
| `Stop()` | SIGINT → V8이 커버리지 JSON 덤프 |
| `Collect()` | `npx c8 report --reporter=json` → 파싱 |
| `Reset()` | v8cov 디렉토리 초기화 |

### c8 JSON 포맷

```json
{
  "result": [{
    "url": "file:///app/routes/user.js",
    "functions": [{
      "functionName": "createUser",
      "ranges": [{
        "startOffset": 100,
        "endOffset": 500,
        "count": 1
      }, {
        "startOffset": 200,
        "endOffset": 300,
        "count": 0
      }]
    }]
  }]
}
```

V8 커버리지는 바이트 오프셋 기반이다. 라인 번호로 변환 필요.

## 파일 목록

| 파일 | 상태 |
|------|------|
| `internal/adapter/node.go` | 신규 |
| `internal/coverage/node.go` | 신규 — V8/c8 JSON 파서 |

## 완료 기준

- [ ] Express 프로젝트에서 V8 커버리지 계측 서버가 기동/종료된다
- [ ] c8 JSON에서 미커버 라인이 추출된다
- [ ] IMPROVE 출력에 Node.js 핸들러의 미커버 라인이 표시된다
