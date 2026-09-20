# MISTAKES — repeating-errors journal

Rule: once logged here — never repeated. One entry = error + symptom + why + fix.

## Repeating (повторяющиеся)

- none yet

## Observations 2026-09-20 (first rake hits, not counted as repeating yet)

1. `time.After` in consumer select-loop (`cmd/order/main.go`) — new timer every round, timer leak; fix with `time.NewTimer`/`Ticker`.
2. Package names off from dirs: `package pros`, `package gRPS`.
3. Module case mismatch: `go.mod` GRPS vs gen-prefix gRPS.
4. Russian comment in `.proto`.
