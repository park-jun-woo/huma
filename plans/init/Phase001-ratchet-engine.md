# Phase 001: Ratchet Engine ✅ DONE

## 목표

`hurlfill next` 하나로 작동하는 래칫 엔진을 완성한다.
실제 Go Gin 프로젝트에서 end-to-end 동작을 검증한다.

## 전제

- 레거시 서버는 사용자가 직접 띄운다. hurlfill은 서버를 기동하지 않는다.
- 검증기는 `hurl --test`다. 커버리지 계측은 이 Phase에 없다.
- hurl PASS/FAIL이 유일한 판정 기준이다.

## 구현

### 1. hurlfill.yaml

hurlfill이 알아야 하는 것만.

```yaml
base_url: "http://localhost:8080"
hurl_dir: "hurl"
scan:
  lang: go
```

### 2. Scanner

현재 regex 기반 스캐너를 유지한다.
AST 파싱은 이 Phase에서 하지 않는다.
Endpoint 구조체에 핸들러 함수의 파일 위치만 정확히 담으면 충분하다.

### 3. Session

래칫 상태 기계. 상태는 세 가지.

```
TODO → (hurl PASS) → PASS
TODO → (hurl FAIL) → TODO (재시도)
PASS → 되돌리지 않음
```

DONE과 IMPROVE는 Phase 003(커버리지)에서 추가한다.
이 Phase에서는 hurl이 통과하면 PASS, 실패하면 TODO 유지.

### 4. next 커맨드

```
hurlfill next 실행 시:
1. session에서 첫 번째 TODO를 찾는다
2. .hurl 파일 없음 → TODO 출력
3. .hurl 파일 있음 → hurl --test 실행
   - FAIL → 에러 출력
   - PASS → session 갱신, 다음 TODO 출력
4. TODO 없음 → "All endpoints complete!"
```

### 5. 검증용 테스트 서버

`testdata/` 디렉토리에 최소 Go Gin 서버를 만든다.
엔드포인트 5개. hurlfill의 end-to-end 테스트 대상.

```
POST   /api/users        — 생성
GET    /api/users         — 목록
GET    /api/users/:id     — 조회
PUT    /api/users/:id     — 수정
DELETE /api/users/:id     — 삭제
```

## 파일 목록

| 파일 | 상태 |
|------|------|
| `internal/config/config.go` | 신규 |
| `cmd/next.go` | 수정 — config 연동 |
| `cmd/scan.go` | 수정 — config 연동 |
| `testdata/server/` | 신규 — 검증용 Gin 서버 |

## 완료 기준

- [ ] testdata 서버를 띄우고 `hurlfill scan` → `hurlfill next` → hurl 작성 → `hurlfill next` → PASS 흐름이 동작한다
- [ ] 5개 엔드포인트 전부 PASS 후 "All endpoints complete!" 출력
- [ ] 에이전트 죽어도 session.json에서 이어간다
