---
name: huma
description: "Ratchet tool for wall-to-wall Hurl API test generation. Use when writing Hurl tests for SaaS backends, when an endpoint lacks test coverage, or when asked to generate API integration tests. Triggers on keywords: hurl, endpoint test, API test, ratchet, coverage."
metadata:
  author: park-jun-woo
  version: "0.3.0"
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

Requires **Go 1.22+**. No cgo dependency — pure Go. Built on the [reins](https://github.com/park-jun-woo/reins) quest framework — only the deterministic gate locks `PASS`.

## Commands

| Command | Purpose |
|---------|---------|
| `huma scan [openapi] [src-root]` | Seed endpoints. Arg 1 = OpenAPI/JSON/YAML (auto-detects `openapi.yaml` if omitted); arg 2 = source root linking handlers to `file:line` (source-branch oracle) |
| `huma next` | Show the next TODO + authoring prompt (**read-only** — does not advance) |
| `huma submit --key <id> [--in <hurl>\|-]` | Evaluate one endpoint's `.hurl` through the CRI gate (static). `--in` is the hurl path; omit for the conventional path |
| `huma cover` | **Live mode**: bring the server up once, run each endpoint's hurl, measure runtime coverage, gate each |
| `huma cover --generate [--model <b:m>] [--max-items N]` | **Unattended loop**: LLM generates `.hurl`, server measures, CRI gate converges |
| `huma status` | Progress (TODO/PASS/IMPROVE/UNVERIFIED/DONE/SKIPPED) with per-endpoint CRI tier |
| `huma rules` | Print the gate's rule catalog (M/E/H/S/A/C) |
| `huma export` | Emit terminal results as JSONL (emit-once) |

## Workflow

```
1. Create manifest.yaml      — project config + testing block
2. huma scan openapi.yaml .  — seed endpoints (arg 2 = source root, optional but enables the oracle)
3. Manual loop:
   a. huma next                — read the TODO prompt (read-only; handler source + responses + hurl example)
   b. Write the .hurl file at the path shown
   c. huma submit --key <id>   — evaluate through the CRI gate → PASS / IMPROVE / UNVERIFIED
   d. Repeat until "All complete"
   OR unattended (live): huma cover --generate --model claude:sonnet
```

> **`next` is read-only.** Unlike older versions, `huma next` only *shows* the prompt — use `huma submit` to evaluate, or `huma cover` to run the live loop over all endpoints.

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
huma scan openapi.yaml .       # arg 2 = source root (optional; enables the source-branch oracle)
# seeded 65 item(s); 65 total
```

huma accepts OpenAPI yaml (auto-detected by the `openapi:` key), JSON arrays, or YAML endpoint lists.
Each endpoint becomes a `quest.Item` in `session.json`. Run `huma status` for the tally.

### Step 3: Ratchet loop

```bash
huma next
# # TODO  POST /api/v1/auth/login
# # Source: internal/api/auth/handler.go:412
# # Handler: Login
# ...hurl example + response fields...

#  ...write hurl/post_api_v1_auth_login.hurl...

huma submit --key <id>         # evaluate → PASS / IMPROVE / UNVERIFIED
```

`huma next` is read-only (shows the prompt only). After writing the `.hurl`, run `huma submit --key <id>` to evaluate it; on PASS the ratchet locks it. For full live verification across all endpoints, use `huma cover` (or `huma cover --generate` to let an LLM author and converge them unattended).

## Key Concepts

### Ratchet States

| State | Meaning |
|-------|---------|
| **TODO** | No .hurl file yet |
| **PASS** | All expected client branches covered, at or above the required evidence tier (CRI) |
| **IMPROVE** | Hurl passes but coverage < 100%; uncovered lines/branches shown |
| **UNVERIFIED** | No independent oracle — source unlinked **and** server uninstrumented. Measurement failed, so **not** a pass. Fix by passing a source root to `huma scan` or adding `testing.server` |
| **DONE** | Coverage stalled after retry **and** every uncovered branch has a verifiable reason in `.huma/unreachable.yaml`; accepted at current level |

### Evidence tiers (CRI)

A `PASS` is not a single state — it carries a **cheese-resistance index** (CRI 0–3) so a weak verdict cannot masquerade as a strong one. `huma status` prints the tier per endpoint.

| Label | CRI | Meaning |
|-------|-----|---------|
| **UNVERIFIED** | 0 | No oracle / no execution / no denominator. Never a pass. |
| **SCAFFOLDED** | 1 | Hurl written, not executed (static-mode honest ceiling). |
| **SMOKE** | 2 | Server executed green; no branch-line binding (uninstrumented). |
| **COVERED** | 3 | Source ∪ runtime: every client branch is runtime-bound and asserted. |

Set the minimum tier for `PASS` via `testing.require_cri` in `manifest.yaml`. If unset, huma auto-requires the maximum tier reachable in the current mode. **Key rule: no-signal is never a pass** — a thin OpenAPI (one status per endpoint), an unlinked source, or an uninstrumented server yields UNVERIFIED, not a free 100%.

### Hurl File Naming

huma expects: `{hurl_dir}/{method}_{path_with_underscores}.hurl`

Example: `GET /api/v1/admin/buildings` → `hurl/get_api_v1_admin_buildings.hurl`

### Static Mode (no server)

Without `testing.server`, huma uses OpenAPI responses and/or source code static analysis to:
- Show expected response status codes in TODO prompt
- Check hurl files for status code coverage without running them
- Mark PASS (tier SCAFFOLDED) when all expected codes are covered in hurl assertions

**Static mode does not execute hurl** — it verifies the file is *written* to cover the codes, not that it *passes*. If there is no source link and no server, there is no independent oracle, so the verdict is **UNVERIFIED**, not PASS. Provide a source root as the second `scan` arg (`huma scan openapi.yaml .`) so the denominator comes from real handler branches (ground truth) — OpenAPI declarations can only *add* branches, never shrink them.

### Live Mode (`huma cover`)

With `testing.server` config, `huma cover` brings the server up once and, per endpoint:
1. Builds the instrumented binary (once)
2. Resets coverage → starts the server → waits for the ready endpoint
3. Runs the endpoint's hurl test
4. Stops the server (Go flushes coverage only on exit) → collects per-handler line coverage
5. Feeds runtime coverage into the CRI gate → PASS (SMOKE/COVERED) or IMPROVE with uncovered lines

Add `--generate` to have an LLM author each `.hurl` and converge it through the gate unattended. (Go targets must handle `SIGINT` with a graceful shutdown so coverage counters flush.)

## Common Errors and Fixes

All errors carry a rule ID. See `rulebook.md` for the full catalog.

| Rule ID | Error | Cause | Fix |
|---------|-------|-------|-----|
| S-01 | `[S-01] No session found` | Haven't scanned yet | Run `huma scan openapi.yaml` |
| H-01 | `[H-01] Hurl file not found at expected path` | .hurl not at expected path | Check `huma next` output for expected filename |
| M-02 | `[M-02] manifest.yaml parse error` | Bad manifest.yaml | Validate YAML syntax |
| E-01 | `[E-01] No OpenAPI file found` | No openapi.yaml and no path arg given | Place openapi.yaml in project root or pass it: `huma scan <path>` |
| H-02 | `[H-02] Hurl execution failed` | hurl binary not installed or server not running | Install hurl: `cargo install hurl` or `brew install hurl` |
| A-02 | `[A-02] Server build command failed` | Build command errored | Check build command in manifest.yaml testing.server.build |
| C-01 | Verdict downgraded to UNVERIFIED | No-signal cannot PASS (no oracle/execution/denominator) | Pass a source root to `huma scan`, or add `testing.server`, so coverage can be measured |
| C-02 | Denominator must be monotonic | Input spec tried to shrink ground-truth branches | Don't rely on a thin OpenAPI; source branches set the floor |
| C-03 | Assertion depth too shallow | Entry asserts status only (or is skip/empty) | Add body/header assertions for the branch |
| C-04 | DONE requires a reason | Uncovered branch has no justification | Record it in `.huma/unreachable.yaml` (reason + source evidence) |

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
| Fastify | `fastify` | NodeAdapter | regex |
| Hono | `hono` | NodeAdapter | regex |
| Supabase Edge Functions | `deno` | DenoAdapter | regex |

## Full Documentation

| Document | Purpose |
|----------|---------|
| `SKILL.md` | Agent workflow and commands |
| `rulebook.md` | All 37 validation rules with IDs |
| `README.md` | Quick start and overview |
