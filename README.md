# gRPS

Учебный боевой проект: gRPC-сервис заказов, доведённый до продакшн-уровня.
Дедлайн: апрель 2027 — уверенное прохождение техсобеса уровня middle по Go-бэкенду.
Стек по мере ввода тем: Go → gRPC → Postgres → Docker → Redis → Kafka → Kubernetes → observability → CI/CD.

## Раскладка папок

| Путь | Что это | В git? |
|---|---|---|
| `proto/orders/v1/order.proto` | исходник контракта (single source of truth — единственный источник правды) | ДА |
| `gen/orders/v1/` | сгенерированный код, только здесь и нигде больше | НЕТ (`.gitignore`: `*.pb.go`) |
| `main.go` | точка входа | ДА |
| `orders/` в корне репо | мёртвый мусор от старой конфигурации, подлежит удалению | должен исчезнуть |

## Формула путей

```
proto/orders/v1/order.proto
        │ buf generate
        ▼
gen/orders/v1/order.pb.go        ← типы (messages)
gen/orders/v1/order_grpc.pb.go   ← клиент/сервер gRPC (stubs — заготовки вызовов)
```

- Import path в Go: `github.com/RastBast/gRPS/gen/orders/v1`.
- `option go_package` в .proto руками НЕ пишем — его подставляет managed mode из `buf.gen.yaml` (go_package_prefix = `github.com/RastBast/gRPS/gen`).
- Два `.go` файла после генерации — норма, а не баг: два плагина, две работы.

## Порядок команд (жёсткий)

1. `buf generate`
2. `go mod tidy`
3. `go build ./...`

## Проверка, что пуш реально долетел

```bash
git rev-parse HEAD          # хэш твоего локального коммита
git ls-remote origin master # хэш на GitHub
# совпали → пуш существует. Не совпали → пуша нет, как бы ты ни был уверен.
```

## Файлы памяти (читать в начале каждого нового чата)

- `PROGRESS.md` — чек-лист тем, квизы, лог мотивации M/E.
- `TECH_DEBT.md` — код, написанный ИИ: обязан переписать и объяснить.
- `MISTAKES.md` — журнал повторяющихся ошибок, перечитывать перед коммитом.
