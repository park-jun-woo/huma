---
name: huma
description: "Auto-generate wall-to-wall Hurl tests for legacy SaaS APIs"
---

Run `huma next` and follow the instructions in the output.

Rules:
- The legacy server's current behavior is the oracle.
  Write assertions matching what the server ACTUALLY returns, not what it SHOULD return.
- If hurl test fails, read the error, fix the .hurl file, run `huma next` again.
- If huma next says "All endpoints complete!", stop.
- Do NOT skip endpoints. Do NOT declare completion yourself.
  Only `huma next` decides when you're done.

Repeat until "All endpoints complete!" appears.
