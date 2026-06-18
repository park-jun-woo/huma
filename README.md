# huma

<p align="center">
  <img src="huma.webp" alt="huma" width="480" />
</p>

[![Version](https://img.shields.io/badge/version-v0.3.0-blue.svg)](https://github.com/park-jun-woo/huma/releases)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![skills.sh](https://skills.sh/b/park-jun-woo/huma)](https://skills.sh/park-jun-woo/huma)

> **Recommended:** [Claude Code](https://claude.ai/code). Tested and optimized for Claude Code.

**Hurl Master** — Ratchet tool that drives AI agents to write wall-to-wall [Hurl](https://hurl.dev) tests for SaaS backend APIs.

## Why huma

- **0 → 100% endpoint coverage** — Point huma at any backend. Every endpoint gets a Hurl test. No exceptions.
- **Agent-native** — `huma next` outputs exactly what an AI agent needs: handler source, expected responses, hurl example, file path. Copy-paste-free loop.
- **Ratchet guarantee** — Coverage only goes up, never down. Once an endpoint passes, it stays passed.
- **OpenAPI-first** — Feed it an OpenAPI yaml and go. No manual endpoint listing needed.
- **Two modes** — Static mode (no server, analysis only) for scaffolding. Live mode (server running) for full verification.
- **Language-agnostic** — Go, Python, Node.js. Same workflow, same ratchet.

## Quick start

```bash
npx skills add park-jun-woo/huma
```

That's it. Your AI agent now knows how to run the full ratchet loop. See [SKILL.md](SKILL.md) for details.

## How it works

huma is built on the [reins](https://github.com/park-jun-woo/reins) quest framework: **generation is probabilistic, the gate is deterministic, and only the gate locks `PASS`.**

```
openapi.yaml ──► huma scan ──► session
                                  │
                            huma next ──► (read-only) TODO prompt
                              │            handler source + responses + hurl example
                        agent writes .hurl
                              │
                       huma submit ──► CRI gate ──► PASS / IMPROVE / UNVERIFIED
                              ▲                          │
                              └──────── retry ◄──────────┘

   unattended:  huma loop                  (LLM writes .hurl, server measures runtime coverage, gate converges)
                huma loop --measure-only    (measure existing .hurl, no LLM)
```

```bash
# 1. Seed endpoints from OpenAPI (+ optional source-link root for the branch oracle)
huma scan openapi.yaml .          # or just `huma scan` to auto-detect openapi.yaml

# 2a. Manual loop
huma next                          # show next TODO + authoring prompt (read-only)
#    ...agent writes the .hurl file...
huma submit --key <id>             # evaluate that endpoint through the CRI gate
huma status                        # progress + per-endpoint CRI tier

# 2b. Or unattended: the server runs, an LLM generates each endpoint's .hurl, and runtime
#     coverage is measured and converged through the gate
huma loop                          # --model defaults to ollama:gemma4:e4b
huma loop --model claude:sonnet --max-items 50
huma loop --measure-only           # measure existing .hurl only; no LLM
```

## Commands

| Command | Description |
|---------|-------------|
| `huma scan [openapi] [src-root]` | Seed endpoints. Arg 1 = OpenAPI/JSON/YAML (auto-detects `openapi.yaml` if omitted); arg 2 = source root to link handlers to `file:line` (enables the source-branch oracle) |
| `huma next` | Show the next TODO + authoring prompt (**read-only** — does not advance) |
| `huma submit --key <id> [--in <hurl>\|-]` | Evaluate one endpoint's `.hurl` through the CRI gate (static). `--in` is the hurl path; omit to use the conventional path |
| `huma loop [--model <b:m>] [--max-items N]` | **Unattended live loop**: bring the server up once, have an LLM generate each remaining TODO's `.hurl`, run it, measure runtime coverage, and converge through the CRI gate to PASS/DONE. `--model` defaults to `ollama:gemma4:e4b` |
| `huma loop --measure-only` | Measure existing `.hurl` against the live server with no LLM generation — coverage measurement only |
| `huma status` | Progress tally (TODO/PASS/IMPROVE/UNVERIFIED/DONE/SKIPPED) with per-endpoint CRI tier |
| `huma rules` | Print the gate's rule catalog (M/E/H/S/A/C) |
| `huma export` | Emit terminal results as JSONL (emit-once) |

## Two modes

### Static mode (no server)

When `manifest.yaml` has no `testing.server` block, huma works without a running server:

- Extracts expected response status codes from OpenAPI or source analysis
- Checks hurl files for status code coverage
- Guides the agent to write hurl test scaffolds

> **Static mode does not execute hurl.** It verifies that a `.hurl` file is *written* to cover the expected status codes — not that it *passes* when run. A scaffold with a wrong URL, bad assertion, or syntax error can still reach `PASS`. See [Scope](#scope-what-huma-guarantees) below.

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

### Live mode (server running)

Add `testing.server` to run hurl tests against a live server with runtime coverage:

```yaml
testing:
  # ... same as above, plus:
  server:
    build: "go build -o ./server.test ./cmd/server"
    start: "./server.test"
    ready: "/api/health"
```

`huma loop` then builds the instrumented server, brings it up once, and per endpoint resets coverage → runs the hurl → stops (to flush coverage) → collects. By default an LLM authors each `.hurl` and converges it through the gate unattended; pass `--measure-only` to measure existing `.hurl` without generation.

> **Go coverage note:** Go flushes integration coverage only on process exit, so `loop` cycles the server per endpoint and the target server must handle `SIGINT` with a graceful shutdown for counters to be written.

### Auth & fixtures: dynamic `{{token}}` (live mode)

Protected endpoints need an `Authorization: Bearer {{token}}` value that static `hurl_variables` can't supply (tokens expire, fixture IDs are created at runtime). `huma loop` resolves these **once before the loop** and injects them into **every** endpoint's hurl run (both generated and `--measure-only`). Captured/minted variables override `hurl_variables` on a name clash.

**Capture (recommended)** — `testing.setup` points at a user-authored `.hurl` that logs in and captures the token (and any fixtures) via hurl's native `[Captures]`:

```yaml
testing:
  # ... server config as above, plus:
  setup:
    hurl: "setup/auth.hurl"
```

`setup/auth.hurl`:

```
POST {{host}}/api/v1/auth/login
Content-Type: application/json
{"email":"admin@example.com","password":"secret"}
HTTP 200
[Captures]
token: jsonpath "$.token"
# fixtures come free — capture anything the setup hurl declares:
# building_id: jsonpath "$.building.id"
```

Every variable the setup hurl captures (`token`, `building_id`, …) is injected as `{{token}}`, `{{building_id}}`, etc. Capture happens once per loop — a stateless JWT stays valid across the per-endpoint server restarts.

**Mint (option, no login)** — when you hold the signing secret, `testing.auth` hand-signs an HS256 token directly, skipping the login endpoint:

```yaml
testing:
  auth:
    type: "jwt-hs256"
    secret_env: "GOZHIP_JWT_SECRET"   # secret read from this env var
    claims: { role: "admin", sub: "1", twofa_complete: "true" }
```

huma signs `{header,claims+exp}` with the env secret and injects the result as `{{token}}`. Limitation: the claim value types and algorithm must match what your app expects (claims are emitted as strings; RS256/ES256 are out of scope — use capture for those).

> **Prerequisite — seed an admin user:** the capture path logs in, so a login-capable admin user must already exist in the DB. User seeding is app-specific (schema/hash dependent), so huma does **not** do it: put your seed command in `testing.deps.up` or seed before running. The mint path bypasses seeding (signature-only). If capture/mint fails, huma warns loudly and continues token-less so you still see which endpoints 401.

### NestJS + Supabase

```yaml
apiVersion: yongol/v1
kind: Project
metadata:
  name: my-nestjs-app
backend:
  lang: node
  framework: nestjs
  module: my-nestjs-app
testing:
  base_url: "http://localhost:3000"
  hurl_dir: "hurl"
  hurl_variables:
    host: "http://localhost:3000"
    supabase_url: "http://127.0.0.1:54321"
    supabase_anon_key: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    supabase_service_role_key: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  deps:
    up: "supabase start"
    down: "supabase stop"
    ready: "http://127.0.0.1:54321/rest/v1/"
  server:
    build: "npm run build"
    start: "node dist/main.js"
    ready: "/health"
```

### Express

```yaml
apiVersion: yongol/v1
kind: Project
metadata:
  name: my-express-app
backend:
  lang: express
  framework: express
  module: my-express-app
testing:
  base_url: "http://localhost:3000"
  hurl_dir: "hurl"
  hurl_variables:
    host: "http://localhost:3000"
  server:
    build: "npm install"
    start: "node src/index.js"
    ready: "/health"
```

### Supabase Edge Functions

```yaml
apiVersion: yongol/v1
kind: Project
metadata:
  name: my-edge-functions
backend:
  lang: deno
  framework: supabase
  module: my-edge-functions
testing:
  base_url: "http://localhost:54321"
  hurl_dir: "hurl"
  hurl_variables:
    host: "http://localhost:54321"
    supabase_anon_key: "eyJ..."
    supabase_service_role_key: "eyJ..."
  deps:
    up: "supabase start"
    down: "supabase stop"
  server:
    start: "supabase functions serve"
    ready: "/functions/v1/health"
```

## Ratchet states

| State | Meaning |
|-------|---------|
| **TODO** | No .hurl file. Agent gets handler source + expected responses + hurl example. |
| **IMPROVE** | Hurl exists but missing response status codes. Agent gets specific missing codes. |
| **UNVERIFIED** | No independent oracle: source is unlinked **and** the server is uninstrumented (no runtime evidence). Measurement failed, so this is **not** a pass. Fix by passing a source root to `huma scan` or adding `testing.server`. |
| **PASS** | All expected client branches covered, at or above the required evidence tier (CRI). |
| **DONE** | Coverage stalled after retry, **and** every uncovered branch has a verifiable reason in `.huma/unreachable.yaml`. Accepted at current level. |

### Evidence tiers (CRI)

A `PASS` is not a single state — it carries a **cheese-resistance index** (CRI 0–3) so a weak verdict cannot masquerade as a strong one. `huma status` prints the tier for every endpoint.

| Label | CRI | Meaning |
|-------|-----|---------|
| **UNVERIFIED** | 0 | No oracle / no execution / no denominator. Never a pass. |
| **SCAFFOLDED** | 1 | Hurl written, not executed (static mode honest ceiling). |
| **SMOKE** | 2 | Server executed green; no branch-line binding (uninstrumented). |
| **COVERED** | 3 | Source ∪ runtime: every client branch is runtime-bound and asserted. |

The minimum tier required for `PASS` is set by `testing.require_cri` in `manifest.yaml`. If unset, huma auto-requires the maximum tier reachable in the current mode.

## Scope: what huma guarantees

huma is a **coverage ratchet**, not a test runner you can trust blindly. Be clear about what each mode verifies:

| | huma guarantees | huma does **not** guarantee |
|---|---|---|
| **Static mode** | A `.hurl` file exists and is written to cover every expected response status code | That the hurl test actually **passes** — the server is never contacted |
| **Live mode** | The hurl test was executed and exited green against a running server, with runtime coverage measured | That your assertions are *meaningful* — a test can pass while asserting too little |

**Always run a hurl smoke test yourself before trusting coverage.** A green `PASS` from huma means "covered," not "verified" — especially in static mode, where a scaffold with a wrong URL, broken assertion, or syntax error never gets executed.

```bash
# Smoke-test everything huma generated against a live server
hurl --test --variable host=http://localhost:8080 hurl/*.hurl
```

If the smoke test fails, the `.hurl` files need fixing regardless of what huma reports. Treat huma's ratchet as "no endpoint left *unwritten*"; treat the smoke test as "no endpoint left *unverified*."

### Residual gaps even at COVERED (CRI 3)

`COVERED` means "every client branch is runtime-bound and asserted," **not** "semantically correct." These gaps are out of scope by design:

- **R-1 Semantic correctness** — A response whose shape matches but whose *values* are wrong (schema passes, business meaning violated). huma does not interpret meaning; this is the job of SSOT/domain assertions (e.g. yongol).
- **R-2 Analyzer false-negatives** — regex-based language analyzers (node/python/php/java/rust/deno) can miss branches (dynamic routing, etc.), so the denominator may be incomplete. The Go AST analyzer is the most complete; for regex languages treat the branch list as a floor, not a ceiling.
- **R-3 Oracle corruption** — A miswired instrumented server that always returns 200 makes the execution evidence (E axis) lie. A ready-check alone cannot detect this; pair it with a deliberate-error sanity probe.

## Supported languages

| Language | Adapter | Analyzer | `backend.lang` |
|----------|---------|----------|-----------------|
| Go | GoAdapter | go/ast | `go` (default) |
| Python | PythonAdapter | regex | `python` |
| Node.js | NodeAdapter | regex | `node` |
| NestJS | NodeAdapter | regex | `nestjs` |
| Express | NodeAdapter | regex | `express` |
| Fastify | NodeAdapter | regex | `fastify` |
| Hono | NodeAdapter | regex | `hono` |
| Supabase Edge Functions | DenoAdapter | regex | `deno` |

## Pipeline

huma fits in the juicer → huma → yongol pipeline:

```
Legacy codebase
    │
    ▼
juicer ──► openapi.yaml    (extract API spec)
    │
    ▼
huma ──► hurl/*.hurl        (generate tests)
    │
    ▼
yongol ──► refactored code  (SSOT-based rebuild)
```

`manifest.yaml` is shared across huma and yongol — zero transition cost.

## Install

```bash
go install github.com/park-jun-woo/huma@latest
```

Or clone and build:

```bash
git clone https://github.com/park-jun-woo/huma.git
cd huma && make install
```

## License

MIT — see [LICENSE](LICENSE).
