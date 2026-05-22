# Phase 002: Harness — Prompt + Skill + Loop Driver

## 목표

hurlfill을 CLI 도구에서 하네스 엔지니어링 스킬로 만든다.
Phase 001의 래칫 엔진 위에 세 가지를 올린다:
1. 에이전트가 따를 수 있는 출력 (prompt generation)
2. Claude Code 스킬 정의 (loop driver)
3. 비-Claude 에이전트를 위한 단건 프롬프트 생성 (prompt export)

## 핵심 원칙: 오라클로서의 레거시 코드

레거시 서버의 현재 동작이 진실이다.
버그가 있어도 그게 오라클이다. 테스트는 현재 동작을 캡처한다.

```
목적: "이 API가 올바르게 동작하는가?" ← 아님
목적: "이 API가 지금과 동일하게 동작하는가?" ← 이것
```

에이전트에게 이 관점을 명시적으로 전달한다.

## 1. 출력 = 에이전트 지시문

`hurlfill next`의 stdout이 완전한 에이전트 프롬프트가 된다.
IFEval 수렴 조건 세 가지를 충족시킨다.

### 조건 1: 결정론적 사실

```
# TODO  POST /api/users
# Source: internal/handler/user.go:30
# Handler: CreateUser
```

### 조건 2: 예시가 컨텍스트에 존재

```
## Handler source

  func (h *Handler) CreateUser(c *gin.Context) {
      var req CreateUserRequest
      if err := c.ShouldBindJSON(&req); err != nil {
          c.JSON(400, gin.H{"error": err.Error()})
          return
      }
      ...
  }

## Hurl example

  POST http://localhost:8080/api/users
  Content-Type: application/json
  {"name": "test", "email": "test@example.com"}

  HTTP 200
  [Asserts]
  jsonpath "$.id" exists
```

handler source를 직접 읽어서 출력에 포함한다.
에이전트가 파일을 찾아 읽을 필요가 없다.

### 조건 3: 단건 지시 (broad exploration 차단)

```
## Instructions

1. Write hurl/post_api_users.hurl
2. Run `hurlfill next`
```

한 번에 하나. Read, Write, Run. 세 단어.

## 2. Claude Code 스킬 — 루프 드라이버

스킬이 루프를 돌린다. 에이전트는 루프의 존재를 모른다.

```yaml
# .claude/skills/hurlfill.md
---
name: hurlfill
description: "Auto-generate wall-to-wall Hurl tests for legacy SaaS APIs"
---

Run `hurlfill next` and follow the instructions in the output.

Rules:
- The legacy server's current behavior is the oracle.
  Write assertions matching what the server ACTUALLY returns, not what it SHOULD return.
- If hurl test fails, read the error, fix the .hurl file, run `hurlfill next` again.
- If hurlfill next says "All endpoints complete!", stop.
- Do NOT skip endpoints. Do NOT declare completion yourself.
  Only `hurlfill next` decides when you're done.

Repeat until "All endpoints complete!" appears.
```

Claude는 이 스킬을 따른다. IFEval이 높은 모델이기 때문이다.

## 3. 비-Claude 에이전트 — 단건 프롬프트 생성

Codex, Grok Build 등은 루프를 안 돌린다.
이들에게는 루프를 요구하지 않는다. 단건 태스크를 던진다.

```bash
hurlfill prompt
```

현재 TODO 엔드포인트에 대한 완전한 단건 프롬프트를 stdout으로 출력한다.
handler source + hurl example + instruction이 전부 포함된, 자기완결적 프롬프트.

외부 루프는 hurlfill 바깥에서 돈다:

```bash
while hurlfill prompt > /tmp/task.txt 2>/dev/null; do
    cat /tmp/task.txt | codex --prompt -
    hurlfill verify            # hurl --test 실행 + session 갱신
done
```

hurlfill은 프롬프트 생성과 검증만 한다. 에이전트 호출은 바깥의 몫이다.
루프 드라이버 문제를 hurlfill 안에서 풀지 않는다.

## 파일 목록

| 파일 | 상태 |
|------|------|
| `cmd/next.go` | 수정 — handler source + example 포함 출력 |
| `cmd/prompt.go` | 신규 — 단건 프롬프트 생성 |
| `cmd/verify.go` | 신규 — hurl 실행 + session 갱신 (prompt 모드용) |
| `skill/hurlfill.md` | 신규 — Claude Code 스킬 정의 |

## 완료 기준

- [ ] `hurlfill next` 출력에 handler source와 hurl example이 포함된다
- [ ] Claude Code 스킬로 testdata 서버의 5개 엔드포인트를 전부 커버한다
- [ ] `hurlfill prompt` + 외부 루프로도 동일하게 동작한다
- [ ] 에이전트가 오라클 관점으로 테스트를 작성한다 (현재 동작 캡처)
