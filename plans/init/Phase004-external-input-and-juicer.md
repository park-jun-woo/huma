# Phase 004: External Endpoint Input + juicer Integration

## 목표

hurlfill의 자체 스캐너(regex)를 제거하고, 외부 도구의 출력을 표준 입력으로 받는다.
엔드포인트 추출은 hurlfill의 책임이 아니다. 래칫을 돌리는 것이 책임이다.

## 배경

프레임워크마다 라우트 추출 방법이 다르다:

| 프레임워크 | 추출 도구 |
|-----------|----------|
| Go Gin | `juicer scan -json` |
| Django | `manage.py show_urls --format json` |
| FastAPI | `/openapi.json` |
| Express | `express-list-endpoints` |
| Spring Boot | `/actuator/mappings` |
| Rails | `rake routes` |

hurlfill이 각각을 구현하면 끝이 없다.
표준 입력 포맷만 정의하면 어떤 추출기든 붙일 수 있다.

## 표준 엔드포인트 포맷

```json
[
  {
    "method": "POST",
    "path": "/api/v1/auth/login",
    "handler": "Login",
    "file": "internal/api/auth/handler.go",
    "line": 42
  },
  {
    "method": "GET",
    "path": "/api/v1/admin/buildings",
    "handler": "ListBuildings",
    "file": "internal/api/building/handler.go",
    "line": 15
  }
]
```

필수: `method`, `path`
선택: `handler`, `file`, `line` (있으면 handler source 출력에 사용)

## 구현

### 1. `hurlfill scan --from`

```bash
# 파일에서
hurlfill scan --from endpoints.json

# stdin에서 (파이프)
juicer scan ./project -json | hurlfill scan --from -

# 기존 방식 제거
hurlfill scan        # → 에러: "--from 필수"
```

### 2. hurl 변수 설정

gozhip의 기존 hurl 파일은 `{{host}}`를 쓴다. hurlfill은 `{{base_url}}`로 하드코딩되어 있다.
hurlfill.yaml에서 변수명을 설정 가능하게 한다.

```yaml
base_url: "http://localhost:9881"
hurl_dir: "hurl"
hurl_variables:
  host: "http://localhost:9881"
  token: ""
  buildingId: "1"
```

runner가 모든 변수를 `--variable key=value`로 전달한다.

### 3. internal/scanner/ 제거

자체 regex 스캐너를 제거한다. 외부 입력으로 완전 대체.
`internal/scanner/scanner.go`의 `Endpoint` 구조체는 유지하되 `Scan()` 함수를 제거하고
`ParseEndpoints(data []byte) ([]Endpoint, error)` 로 교체한다.

### 4. juicer 출력 → hurlfill 입력 변환

juicer의 출력 포맷:

```yaml
endpoints:
  - method: GET
    path: /api/v1/admin/buildings
    handler: internal/api/building/handler.go:ListBuildings
    file: internal/api/building/handler.go
    line: 15
```

`handler` 필드가 `file:funcName` 형식이다. 파싱해서 Endpoint 구조체에 매핑한다.

## 검증

gozhip admin 백엔드로 end-to-end 테스트:

```bash
cd ~/.clari/repos/gozhip/artifacts/backend/admin
juicer scan . -json | hurlfill scan --from -
hurlfill status    # 108개 엔드포인트, 102개 기존 hurl 매칭 확인
hurlfill next      # 남은 6개에 대해 래칫 시작
```

## 파일 목록

| 파일 | 상태 |
|------|------|
| `cmd/scan.go` | 수정 — `--from` 플래그, 파일/stdin 파싱 |
| `internal/scanner/scanner.go` | 수정 — `Scan()` 제거, `ParseEndpoints()` 추가 |
| `internal/config/config.go` | 수정 — `HurlVariables map[string]string` 추가 |
| `internal/runner/runner.go` | 수정 — 변수 목록을 `--variable`로 전달 |
| `internal/prompt/prompt.go` | 수정 — 변수명을 config에서 읽어 예시에 반영 |

## 완료 기준

- [ ] `hurlfill scan --from endpoints.json`으로 session이 생성된다
- [ ] `juicer scan -json | hurlfill scan --from -` 파이프가 작동한다
- [ ] hurlfill.yaml의 hurl_variables가 runner에 전달된다
- [ ] gozhip에서 `hurlfill status` 실행 시 엔드포인트 수와 기존 hurl 매칭이 정확하다
