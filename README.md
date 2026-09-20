# gRPS — OrderService (Go + gRPC)

Study-to-production project. Goal: middle Go backend by April 2027.
Docs in Russian, code/names/commits/architecture in English.

## Layout (раскладка папок)

```text
proto/orders/v1/order.proto   # contract: Order + OrderService (CreateOrder, GetOrder, GetOrders)
gen/                          # generated: NEVER COMMIT, produced by buf generate
cmd/order/main.go             # entrypoint: producer-consumer demo (NOT a gRPC server yet)
internal/order/domain/        # domain: Order, Status, Reporter
internal/order/usecase/       # use cases: business rules (empty for now)
internal/presentation/grpc/   # transport: gRPC server layer (empty for now)
```

## Path formula (формула путей)

```text
proto/orders/v1/order.proto
  --[buf generate: protoc-gen-go + protoc-gen-go-grpc, paths=source_relative]-->
gen/orders/v1/order.pb.go        # messages: struct Order + marshal/unmarshal
gen/orders/v1/order_grpc.pb.go   # transport: Server interface + Client stub + Register
import path: github.com/RastBast/gRPS/gen/orders/v1
```

No `go_package` in .proto — buf injects it (managed mode). Never by hand.

## Order (порядок, strict)

```text
buf generate → go mod tidy → go build ./...
```

## Commands (команды)

```sh
buf generate          # generate gen/ from proto/
go mod tidy           # sync go.mod/go.sum with imports
go build ./...        # build all packages
go test ./...         # run all tests
find gen -type f      # show what was generated
```

## Memory files (память между чатами)

- `PROGRESS.md` — topic checklist, quizzes, M/E log
- `TECH_DEBT.md` — what was NOT written by the student
- `MISTAKES.md` — repeating-errors journal

## Links (доки)

- buf docs: <https://buf.build/docs>
- gRPC Go quickstart: <https://grpc.io/docs/languages/go/quickstart/>
- Protobuf Style Guide: <https://protobuf.dev/programming-guides/style/>
- Google API Design Guide / AIP: <https://google.aip.dev/>
- Go gRPC package: <https://pkg.go.dev/google.golang.org/grpc>
