# Phase 005: Python Coverage Adapter

## 목표

Python 백엔드(Django, FastAPI, Flask)에 대한 커버리지 어댑터를 구현한다.
라우트 추출은 사용자가 외부 도구로 제공한다 (Phase 004의 --from).
hurlfill은 커버리지 측정만 담당한다.

## 라우트 추출 (사용자 몫)

| 프레임워크 | 방법 |
|-----------|------|
| Django | `manage.py show_urls --format json` |
| FastAPI | `curl localhost:8000/openapi.json` → 변환 |
| Flask | 스크립트: `app.url_map.iter_rules()` |

결과를 hurlfill 표준 포맷(Phase 004)으로 변환하여 `hurlfill scan --from`에 전달.

## 커버리지 어댑터

Python은 프레임워크 무관하게 `coverage.py` 하나로 커버리지를 측정한다.

### hurlfill.yaml

```yaml
base_url: "http://localhost:8000"
hurl_dir: "hurl"
server:
  lang: python
  build: ""
  start: "coverage run --source=app manage.py runserver 0:8000 --noreload"
  ready: "http://localhost:8000/health"
  env:
    DJANGO_SETTINGS_MODULE: "config.settings.test"
```

### PythonAdapter

```go
type PythonAdapter struct {
    cfg      *config.ServerConfig
    baseURL  string
    proc     *exec.Cmd
}
```

| 메서드 | 동작 |
|--------|------|
| `Build()` | no-op (Python은 빌드 없음) |
| `Start()` | `coverage run --source=app manage.py runserver` 실행 |
| `WaitReady()` | ready URL 폴링 (GoAdapter와 동일 로직 재사용) |
| `Stop()` | SIGINT → Python atexit 핸들러가 `.coverage` 파일 덤프 |
| `Collect()` | `coverage json -o .hurlfill/cov.json` 실행 → JSON 파싱 |
| `Reset()` | `.coverage` 파일 삭제 |

### Stop() 동작 원리

coverage.py는 atexit 핸들러로 커버리지 데이터를 저장한다.
SIGINT → Python KeyboardInterrupt → 정상 종료 → atexit → `.coverage` 덤프.
SIGKILL은 atexit가 실행되지 않으므로 사용하지 않는다.

### coverage.py JSON 포맷 (`coverage json` 출력)

```json
{
  "files": {
    "app/views/user.py": {
      "executed_lines": [10, 11, 12, 15, 16],
      "missing_lines": [18, 19, 25, 30],
      "summary": {
        "covered_lines": 5,
        "num_statements": 9,
        "percent_covered": 55.5
      }
    }
  }
}
```

`missing_lines`가 미커버 라인. 핸들러 범위(startLine~endLine)로 필터링.

## 파일 목록

| 파일 | 상태 |
|------|------|
| `internal/adapter/python_adapter.go` | 신규 — PythonAdapter 타입 |
| `internal/adapter/new_python_adapter.go` | 신규 — 생성자 |
| `internal/adapter/python_build.go` | 신규 — no-op Build |
| `internal/adapter/python_start.go` | 신규 — coverage run 실행 |
| `internal/adapter/python_stop.go` | 신규 — SIGINT + .coverage 확인 |
| `internal/adapter/python_collect.go` | 신규 — coverage json 실행 + 파싱 |
| `internal/adapter/python_reset.go` | 신규 — .coverage 삭제 |
| `internal/coverage/parse_coverage_py.go` | 신규 — coverage.py JSON 파서 |
| `cmd/next.go` | 수정 — lang에 따라 어댑터 선택 |

## 완료 기준

- [ ] PythonAdapter가 Adapter 인터페이스를 구현한다
- [ ] coverage.py JSON에서 미커버 라인이 추출된다
- [ ] lang=python일 때 next/verify가 PythonAdapter를 사용한다
- [ ] go test 전부 통과
- [ ] filefunc validate 0 violations
