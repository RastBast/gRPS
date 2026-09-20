# TECH_DEBT — what was NOT written by the student

Rule: only the student writes Go/.proto/Docker. Exception — only on explicit request, logged here, must be rewritten + explained by the student.

| Date       | What | Made by | Must rewrite | Status                                             |
|------------|------|---------|--------------|----------------------------------------------------|
| 2026-09-20 | —    | —       | —            | no debt (student's words; writes Go/proto himself) |

Note: README/PROGRESS/TECH_DEBT/MISTAKES created by mentor on 2026-09-20 — docs/configs (allowed), NOT debt, no rewrite needed.

## Open fix-list (not debt — repair plan, student does it)

- [ ] `go.mod` module case: `GRPS` → `gRPS` (+ import in `cmd/order/main.go`)
- [ ] `package pros` → `usecase`, `package gRPS` → meaningful name
- [ ] root `main.go` stub: delete after diagnostics
- [ ] Russian comment in `order.proto` → English
- [ ] status string → enum; total validation; idempotency; graceful shutdown; pin generator versions
