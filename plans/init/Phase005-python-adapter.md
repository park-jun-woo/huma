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
| `WaitReady()` | ready URL 폴링 |
| `Stop()` | SIGINT → coverage.py가 `.coverage` 파일 덤프 |
| `Collect()` | `coverage json -o .hurlfill/cov.json` → 파싱 |
| `Reset()` | `.coverage` 파일 삭제 |

### coverage.py JSON 포맷

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
| `internal/adapter/python.go` | 신규 |
| `internal/coverage/python.go` | 신규 — coverage.py JSON 파서 |
| `cmd/next.go` | 수정 — lang에 따라 어댑터 선택 |

## 완료 기준

- [ ] Django 프로젝트에서 커버리지 계측 서버가 기동/종료된다
- [ ] coverage.py JSON에서 미커버 라인이 추출된다
- [ ] IMPROVE 출력에 Python 핸들러의 미커버 라인이 표시된다
