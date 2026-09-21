# PROGRESS — checklist to middle (April 2027)

## In work (в работе)

- [ ] `gRPC-01`: buf generate + read gen/ + fix module-path case (debt accepted, see Debt ledger; closes when ledger paid)
- [ ] `gRPC-02`: OrderService server + Go client + first live call (started 2026-09-21, M9 E9, variant B default)

## Closed topics (закрытые темы)

- [x] `proto/orders/v1/order.proto` v1 — own request/response box per rpc, total int64 kopecks, gen/ gitignored. Can defend.
- [x] producer-consumer demo — goroutine/channel/select/context/WaitGroup/Reporter interface (surface level, no error wrapping, no tests).

## Open (не закрыто, carried over)

- `status` string → enum; `total` validation; CreateOrder idempotency; gRPC graceful shutdown; pin protoc-gen-go/protoc-gen-go-grpc versions.

## Debt ledger (принято 2026-09-21, due в gRPC-02)

- Q1 plugin mapping + Q2-2nd-half retell — due with first server run
- Q3 — closes with first live client→server call
- Q4 — closes at first forgotten buf generate
- find gen + build+exit + diff stat — due at first red build in gRPC-02 (paste raw)
- md-fetch (variant A) — due end of week

## Roadmap (порядок тем)

Go basics → gRPC → Postgres → Docker → tests/logs → CI/CD → Redis → Kafka → k8s → observability → system design (from October) + algos (1–2/week → 1–2/day) + behavioral interviews.

## Quiz results (результаты квизов)

- none yet

## M/E log (one line per day)

- 2026-09-20: M9 E9 (chat start, motivation window open — pushing gRPC)
- 2026-09-21: M9 E9, обстановка 6→9 (mindmap: mustEmbed found, client/server interfaces merged — corrected; gRPC-02 coding)

## Mock interviews (пробный собес, monthly)

- none yet

## Retro (ретро, monthly)

- none yet
