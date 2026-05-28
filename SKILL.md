---
name: huma
description: Ratchet tool for wall-to-wall Hurl API test generation. Use when writing Hurl tests for SaaS backends, when an endpoint lacks test coverage, or when asked to generate API integration tests. Triggers on keywords: hurl, endpoint test, API test, ratchet, coverage.
metadata:
  author: park-jun-woo
  version: "0.1.4"
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

Requires **Go 1.22+**. No cgo dependency — pure Go.

## Commands

| Command | Purpose |
|---------|---------|
| `huma scan` | Auto-detect `openapi.yaml` and scan endpoints |
| `huma scan --from <file>` | Scan from OpenAPI, JSON array, or YAML |
| `huma next` | Show next TODO, or verify current and advance |
| `huma verify` | Run hurl test and advance if passing |
| `huma status` | Show progress (TODO/PASS/IMPROVE/DONE) |
| `huma prompt` | Output agent prompt for current TODO (no side effects) |

## Workflow

```
1. Create manifest.yaml      — project config + testing block
2. huma scan --from openapi.yaml   (or just `huma scan` to auto-detect)
3. Loop:
   a. huma next             — read the TODO prompt (includes expected responses)
   b. Write the .hurl file at the path shown
   c. huma next             — verify and advance to next endpoint
   d. Repeat until "All complete"
```

### Step 1: manifest.yaml

```yaml
apiVersion: yongol/v1
kind: Project
metadata:
  name: my-project
backend:
  lang: go
  framework: gin
  module: github.com/org/project
testing:
  base_url: "http://localhost:8080"
  hurl_dir: "hurl"
  hurl_variables:
    host: "http://localhost:8080"
```

Without `testing.server` → **static mode** (no server, static analysis only).
With `testing.server` → **live mode** (server running, hurl execution + runtime coverage).

### Step 2: Scan

```bash
huma scan --from openapi.yaml
# Scanned 65 endpoints
#   PASS:    47 (existing hurl, all responses covered)
#   IMPROVE:  8 (existing hurl, missing responses)
#   TODO:    10 (no hurl file)
```

huma accepts OpenAPI yaml (auto-detected by `openapi:` key), JSON arrays, or YAML endpoint lists.
Existing hurl files are pre-checked during scan — no need to run `huma next` for already-complete endpoints.

### Step 3: Ratchet loop

```bash
huma next
# TODO  POST /api/v1/auth/login
# Handler: Login
#
# ## Expected responses (from OpenAPI)
#   200 — OK
#   400 — Bad Request
#   401 — Unauthorized
#
# ## Instructions
# 1. Write hurl/post_api_v1_auth_login.hurl
# 2. Include test entries for status 200, 400, 401
# 3. Run `huma next`
```

Write the .hurl file, then run `huma next` again. On pass it advances to the next endpoint.

## Key Concepts

### Ratchet States

| State | Meaning |
|-------|---------|
| **TODO** | No .hurl file yet |
| **PASS** | Hurl test passes (coverage 100% or no testing.server block) |
| **IMPROVE** | Hurl passes but coverage < 100%; uncovered lines shown |
| **DONE** | Coverage stalled after retry; accepted at current level |

### Hurl File Naming

huma expects: `{hurl_dir}/{method}_{path_with_underscores}.hurl`

Example: `GET /api/v1/admin/buildings` → `hurl/get_api_v1_admin_buildings.hurl`

### Static Mode (no server)

Without `testing.server`, huma uses OpenAPI responses and/or source code static analysis to:
- Show expected response status codes in TODO prompt
- Check hurl files for status code coverage without running them
- Mark PASS when all expected codes are covered in hurl assertions

### Live Mode (server running)

With `testing.server` config, huma:
1. Checks server healthcheck (prompts to start if down)
2. Builds instrumented binary
3. Starts server, waits for ready endpoint
4. Runs hurl test
5. Collects per-handler line coverage
6. Reports uncovered lines in IMPROVE prompt

## Common Errors and Fixes

All errors carry a rule ID. See `rulebook.md` for the full catalog.

| Rule ID | Error | Cause | Fix |
|---------|-------|-------|-----|
| S-01 | `[S-01] No session found` | Haven't scanned yet | Run `huma scan --from openapi.yaml` |
| H-01 | `[H-01] Hurl file not found at expected path` | .hurl not at expected path | Check `huma next` output for expected filename |
| M-02 | `[M-02] manifest.yaml parse error` | Bad manifest.yaml | Validate YAML syntax |
| E-01 | `[E-01] No OpenAPI file found` | No openapi.yaml and --from not specified | Place openapi.yaml in project root or use `--from` |
| H-02 | `[H-02] Hurl execution failed` | hurl binary not installed or server not running | Install hurl: `cargo install hurl` or `brew install hurl` |
| A-02 | `[A-02] Server build command failed` | Build command errored | Check build command in manifest.yaml testing.server.build |

## Conventions

- Hurl files use `{{host}}` variable for base URL
- One .hurl file per endpoint
- Include golden path (happy case) + at least one error case (400/401/404)
- File naming follows `method_path.hurl` convention with underscores replacing slashes

## Supported Languages

| Language | `backend.lang` | Adapter | Analyzer |
|----------|----------------|---------|----------|
| Go | `go` (default) | GoAdapter | go/ast |
| Python | `python` | PythonAdapter | regex |
| Node.js | `node` | NodeAdapter | regex |
| NestJS | `nestjs` | NodeAdapter | regex |
| Express | `express` | NodeAdapter | regex |
| Supabase Edge Functions | `deno` | DenoAdapter | regex |

## Full Documentation

| Document | Purpose |
|----------|---------|
| `SKILL.md` | Agent workflow and commands |
| `rulebook.md` | All 33 validation rules with IDs |
| `README.md` | Quick start and overview |
