# Этап 1: outbox в orders-service

Файлы лежат по тем же путям, что и в проекте: скопируй их поверх orders-service/.

Удалить:
- internal/service/mocks/Producer.go (интерфейса Producer в сервисе больше нет)

Новое:
- migrations/20261029120000_create_outbox_table.sql
- internal/models/outbox.go
- internal/outbox/relay.go, relay_test.go

Изменено:
- main.go, .mockery.yml, internal/config/config.go, internal/models/events.go,
  internal/repository/repo.go, repo_test.go, internal/service/service.go, service_test.go,
  internal/service/mocks/OrderRepository.go, internal/kafka/producer/producer.go

Новые переменные окружения (необязательные):
- OUTBOX_POLL_INTERVAL_MS (по умолчанию 500)
- OUTBOX_BATCH_SIZE (по умолчанию 50)

Новых зависимостей в go.mod нет.
