---
name: huma
description: Ratchet tool for wall-to-wall Hurl API test generation. Use when writing Hurl tests for SaaS backends, when an endpoint lacks test coverage, or when asked to generate API integration tests. Triggers on keywords: hurl, endpoint test, API test, ratchet, coverage.
metadata:
  author: park-jun-woo
  version: "0.1.0"
---

# huma — Ratchet-driven Hurl test generator for SaaS APIs

## When to Use This Skill

- A SaaS backend has endpoints without Hurl tests
- You need wall-to-wall API integration test coverage
- The user asks to generate, write, or improve Hurl tests
- The user mentions "huma", "hurl coverage", or "endpoint ratchet"

## Do NOT Use When

- The project is a library with no HTTP endpoints
- The user wants unit tests (not API-level integration tests)
- The user wants Postman/Bruno/other test formats (huma is Hurl-only)

## Install

```bash
go install github.com/park-jun-woo/huma@latest
```

## Commands

| Command | Purpose |
|---------|---------|
| `huma scan --from <file>` | Read endpoints from JSON/YAML and create a session |
| `huma next` | Show next TODO endpoint, or verify current and advance |
| `huma verify` | Run hurl test for current endpoint and advance if passing |
| `huma status` | Show progress (todo/pass/done counts) |
| `huma prompt` | Output agent prompt for current TODO (no side effects) |

## Workflow

```
1. Create endpoints.yaml    — list every endpoint (method, path, handler, file, line)
2. Create huma.yaml         — base_url, hurl_dir, variables
3. huma scan --from endpoints.yaml
4. Loop:
   a. huma next             — read the TODO prompt
   b. Write the .hurl file at the path shown
   c. huma next             — verify and advance to next endpoint
   d. Repeat until "All complete"
```

### Step 1: endpoints.yaml

```yaml
endpoints:
  - method: GET
    path: /api/health
    handler: HealthCheck
    file: cmd/server/main.go
    line: 42
  - method: POST
    path: /api/v1/auth/login
    handler: Login
    file: internal/api/auth/handler.go
    line: 50
```

Fields: `method`, `path`, `handler`, `file` (source location), `line`.
Accepts JSON array, YAML array, or YAML with `endpoints` key.

### Step 2: huma.yaml

```yaml
base_url: "http://localhost:8080"
hurl_dir: "hurl"
hurl_variables:
  host: "http://localhost:8080"
scan:
  lang: "go"        # go | python | node
```

Optional coverage mode (builds and instruments the server):

```yaml
server:
  build: "go build -cover -o ./server.test ./cmd/server"
  start: "./server.test"
  ready: "http://localhost:8080/api/health"
  env:
    GIN_MODE: test
```

### Step 3: Scan

```bash
huma scan --from endpoints.yaml
# Scanned 42 endpoints
```

### Step 4: Ratchet loop

```bash
huma next
# TODO  GET /api/v1/users
# Source: internal/api/user/handler.go:25
# Handler: ListUsers
#
# ## Handler source
# func (h *Handler) ListUsers(c *gin.Context) { ... }
#
# ## Instructions
# 1. Write hurl/get_api_v1_users.hurl
# 2. Run `huma next`
```

Write the .hurl file, then run `huma next` again. On pass it advances to the next endpoint.

## Key Concepts

### Ratchet States

| State | Meaning |
|-------|---------|
| **TODO** | No .hurl file yet |
| **PASS** | Hurl test passes (coverage 100% or no server config) |
| **IMPROVE** | Hurl passes but coverage < 100%; uncovered lines shown |
| **DONE** | Coverage stalled after retry; accepted at current level |

### Hurl File Naming

huma expects: `{hurl_dir}/{method}_{path_with_underscores}.hurl`

Example: `GET /api/v1/admin/buildings` → `hurl/get_api_v1_admin_buildings.hurl`

### Coverage Mode

With `server` config, huma:
1. Builds instrumented binary
2. Starts server, waits for ready endpoint
3. Runs hurl test
4. Collects per-handler line coverage
5. Reports uncovered lines in IMPROVE prompt

Without `server` config: pass/fail only (no coverage tracking).

## Common Errors and Fixes

All errors carry a rule ID. See `rulebook.md` for the full catalog.

| Rule ID | Error | Cause | Fix |
|---------|-------|-------|-----|
| S-01 | `[S-01] No session found` | Haven't scanned yet | Run `huma scan --from endpoints.yaml` |
| H-01 | `[H-01] Hurl file not found at expected path` | .hurl not at expected path | Check `huma next` output for expected filename |
| M-02 | `[M-02] manifest.yaml parse error` | Bad huma.yaml | Validate YAML syntax |
| E-01 | `[E-01] --from flag required` | Missing --from flag | Run `huma scan --from <file>` |
| H-02 | `[H-02] Hurl execution failed` | hurl binary not installed or server not running | Install hurl: `cargo install hurl` or `brew install hurl` |
| A-02 | `[A-02] Server build command failed` | Build command errored | Check build command in huma.yaml server.build |

## Conventions

- Hurl files use `{{host}}` variable for base URL
- One .hurl file per endpoint
- Include golden path (happy case) + at least one error case (400/401/404)
- File naming follows `method_path.hurl` convention with underscores replacing slashes

## Supported Languages

| Language | `scan.lang` | Adapter |
|----------|-------------|---------|
| Go | `go` (default) | GoAdapter |
| Python | `python` | PythonAdapter |
| Node.js | `node` | NodeAdapter |
