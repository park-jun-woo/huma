# huma — Rulebook

All validation rules emitted by huma. Each rule has a unique ID, level, and description.

## Rule ID Scheme

| Prefix | Domain |
|--------|--------|
| `M-` | manifest.yaml validation |
| `E-` | endpoints.yaml validation |
| `H-` | hurl file validation |
| `S-` | session state validation |
| `A-` | adapter/server validation |

## Levels

| Level | Behavior | Meaning |
|-------|----------|---------|
| `ERROR` | Blocks progress, immediate feedback | Must fix |
| `WARNING` | Continues, prints warning | Review recommended |

## M. Manifest

| Rule ID | Level | Description |
|---------|-------|-------------|
| M-01 | ERROR | manifest.yaml not found |
| M-02 | ERROR | manifest.yaml parse error (invalid YAML) |
| M-03 | ERROR | `apiVersion` missing or unsupported |
| M-04 | ERROR | `metadata.name` missing |
| M-05 | ERROR | `backend.lang` missing |
| M-06 | ERROR | `testing.base_url` missing |
| M-07 | ERROR | `testing.hurl_dir` missing |
| M-08 | ERROR | `testing.server.start` missing |
| M-09 | ERROR | `testing.server.ready` missing |
| M-10 | WARNING | `testing.hurl_variables` empty |

## E. Endpoints

| Rule ID | Level | Description |
|---------|-------|-------------|
| E-01 | ERROR | `--from` flag required |
| E-02 | ERROR | Input file not readable |
| E-03 | ERROR | Input is not valid JSON/YAML |
| E-04 | ERROR | Endpoint missing `method` field |
| E-05 | ERROR | Endpoint missing `path` field |
| E-06 | WARNING | Endpoint missing `handler` field |
| E-07 | WARNING | Endpoint missing `file` field |
| E-08 | WARNING | Duplicate endpoint |

## H. Hurl

| Rule ID | Level | Description |
|---------|-------|-------------|
| H-01 | ERROR | Hurl file not found at expected path |
| H-02 | ERROR | Hurl execution failed |
| H-03 | ERROR | Hurl test failed |
| H-04 | WARNING | Existing hurl file name doesn't match naming convention |
| H-05 | WARNING | Hurl file missing `{{host}}` variable |

## S. Session

| Rule ID | Level | Description |
|---------|-------|-------------|
| S-01 | ERROR | No session found |
| S-02 | ERROR | Session file corrupt |
| S-03 | WARNING | Session has stale entries |

## A. Adapter / Server

| Rule ID | Level | Description |
|---------|-------|-------------|
| A-01 | ERROR | Server healthcheck failed |
| A-02 | ERROR | Server build command failed |
| A-03 | ERROR | Server start command failed |
| A-04 | ERROR | Server ready timeout |
| A-05 | ERROR | Coverage data collection failed |
| A-06 | WARNING | deps.ready check failed |

## Output Format

```
[M-01] manifest.yaml not found
  ▶ Create manifest.yaml in the project root.

[H-04] WARNING — Existing hurl file name doesn't match naming convention
  ▶ Rename to match, or huma will treat these endpoints as TODO.

[A-02] Server build command failed
  ▶ go build -cover: exit status 1
```
