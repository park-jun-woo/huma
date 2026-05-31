# huma

<p align="center">
  <img src="huma.webp" alt="huma" width="480" />
</p>

[![Version](https://img.shields.io/badge/version-v0.2.0-blue.svg)](https://github.com/park-jun-woo/huma/releases)
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

```
openapi.yaml ──► huma scan ──► session
                                  │
                            huma next ◄──┐
                              │          │
                      ┌───────┴───────┐  │
                      │  TODO         │  │
                      │  (no .hurl)   │  │
                      └───────┬───────┘  │
                        agent writes     │
                        .hurl file       │
                              │          │
                      ┌───────┴───────┐  │
                      │  PASS/IMPROVE │──┘
                      └───────────────┘
```

```bash
# 1. Scan endpoints from OpenAPI (or auto-detect openapi.yaml)
huma scan --from openapi.yaml

# 2. Ratchet loop
huma next     # shows TODO → agent writes .hurl → run again
huma status   # check progress
```

## Commands

| Command | Description |
|---------|-------------|
| `huma scan` | Auto-detect `openapi.yaml` and scan endpoints |
| `huma scan --from <file>` | Scan from OpenAPI, JSON array, or YAML |
| `huma next` | Show next TODO, or verify current and advance |
| `huma verify` | Run hurl test and advance if passing |
| `huma status` | Show progress (TODO/PASS/IMPROVE/DONE/UNVERIFIED) with per-endpoint CRI tier |
| `huma scan --link-source <root>` | Map OpenAPI handlers to source `file:line` (enables source-branch oracle) |
| `huma prompt` | Output agent prompt for current TODO (no side effects) |

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

huma builds, starts, and instruments the server. Healthcheck failure prompts the agent to start the server first.

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
| **UNVERIFIED** | No independent oracle: source is unlinked **and** the server is uninstrumented (no runtime evidence). Measurement failed, so this is **not** a pass. Fix with `--link-source` or `testing.server`. |
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
