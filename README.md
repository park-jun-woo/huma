# huma

[![Version](https://img.shields.io/badge/version-v0.1.0-blue.svg)](https://github.com/park-jun-woo/huma/releases)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![skills.sh](https://skills.sh/b/park-jun-woo/huma)](https://skills.sh/park-jun-woo/huma)

> **Recommended:** [Claude Code](https://claude.ai/code). Tested and optimized for Claude Code.

**Hurl Master** — Ratchet tool that drives AI agents to write wall-to-wall [Hurl](https://hurl.dev) tests for SaaS backend APIs.

## Why huma

- **0 → 100% endpoint coverage** — Point huma at any backend. Every endpoint gets a Hurl test. No exceptions.
- **Agent-native** — `huma next` outputs exactly what an AI agent needs: handler source, hurl example, file path. Copy-paste-free loop.
- **Ratchet guarantee** — Coverage only goes up, never down. Once an endpoint passes, it stays passed.
- **10 minutes to first test** — `scan → next → write → verify`. Four commands. No config ceremony.
- **Language-agnostic** — Go, Python, Node.js. Same workflow, same ratchet.

## Quick start

```bash
npx skills add park-jun-woo/huma
```

That's it. Your AI agent now knows how to run the full ratchet loop. See [SKILL.md](SKILL.md) for details.

## Commands

| Command | Description |
|---------|-------------|
| `huma scan --from <file>` | Read endpoints from JSON/YAML file (or `-` for stdin) and create a session |
| `huma next` | Show next untested endpoint, or verify current one and advance |
| `huma verify` | Run hurl test for current endpoint and advance if passing |
| `huma status` | Show progress summary (todo/pass/done counts) |
| `huma prompt` | Output agent prompt for current TODO endpoint (no side effects) |

## Ratchet states

Each endpoint progresses through these states:

- **TODO** — No .hurl file yet. `huma next` outputs a prompt with handler source and a hurl example.
- **PASS** — Hurl test passes (and coverage is 100% or no server config).
- **IMPROVE** — Hurl passes but coverage < 100%. Agent gets uncovered lines to add more test entries.
- **DONE** — Coverage stalled after retry. Endpoint is accepted at current coverage.

## Coverage mode

When `huma.yaml` includes a `server` block, huma builds, starts, and instruments the server to measure per-handler line coverage:

```yaml
server:
  build: "go build -cover -o ./server.test ./cmd/server"
  start: "./server.test"
  ready: "http://localhost:8080/api/health"
  env:
    GIN_MODE: test
```

Without the `server` block, huma runs in no-coverage mode (pass/fail only).

## Supported languages

| Language | Adapter | `scan.lang` |
|----------|---------|-------------|
| Go | `GoAdapter` | `go` (default) |
| Python | `PythonAdapter` | `python` |
| Node.js | `NodeAdapter` | `node` |

## Endpoint input format

JSON array, YAML array, or YAML with `endpoints` key:

```yaml
endpoints:
  - method: GET
    path: /api/v1/users
    handler: ListUsers
    file: internal/api/user/handler.go
    line: 25
```

Fields: `method`, `path`, `handler`, `file` (source location), `line`.

## Install

```bash
go install github.com/park-jun-woo/huma@latest
```

Or clone and build with version:

```bash
git clone https://github.com/park-jun-woo/huma.git
cd huma && make install
```

## License

MIT — see [LICENSE](LICENSE).
