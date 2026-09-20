# PROGRESS — checklist to middle (April 2027)

## In work (в работе)

- [ ] `gRPC-01`: buf generate + read gen/ + fix module-path case (started 2026-09-20; 2026-09-21: 4/5 artifacts explained, Register + Q3/Q4 + fixes + build proof pending; memory files committed to arena/01a0c09e-grps)

## Closed topics (закрытые темы)

- [x] `proto/orders/v1/order.proto` v1 — own request/response box per rpc, total int64 kopecks, gen/ gitignored. Can defend.
- [x] producer-consumer demo — goroutine/channel/select/context/WaitGroup/Reporter interface (surface level, no error wrapping, no tests).

## Open (не закрыто, carried over)

- `status` string → enum; `total` validation; CreateOrder idempotency; gRPC graceful shutdown; pin protoc-gen-go/protoc-gen-go-grpc versions.

## Roadmap (порядок тем)

Go basics → gRPC → Postgres → Docker → tests/logs → CI/CD → Redis → Kafka → k8s → observability → system design (from October) + algos (1–2/week → 1–2/day) + behavioral interviews.

## Quiz results (результаты квизов)

- none yet

## M/E log (one line per day)

- 2026-09-20: M9 E9 (chat start, motivation window open — pushing gRPC)

## Mock interviews (пробный собес, monthly)

- none yet

## Retro (ретро, monthly)

- none yet
