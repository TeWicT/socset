# Локальная разработка

## Требования

- Go 1.22+ (уточнится в `go.mod`)
- Docker + Docker Compose
- Make (опционально) / Taskfile
- curl / httpie для ручных проверок

## Целевая структура репозитория

```text
socset/
  README.md
  docs/
  deploy/
    docker-compose.yml
    obs/                  # otel, grafana provisioning
  services/
    api-gateway/
    auth-service/
    profile-service/
    social-graph-service/
    post-service/
    chat-service/
    notification-service/
    search-service/
  pkg/                    # общий код: auth jwt, otel, outbox, logging
  schemas/
    events/               # JSON schemas / proto событий
    openapi/              # публичные спеки
```

Структура кода появится по мере фаз; документация уже описывает целевое состояние.

## Порядок подъёма (когда появится compose)

```bash
# 1. Инфра
docker compose -f deploy/docker-compose.yml up -d postgres kafka redis minio elasticsearch jaeger

# 2. Миграции сервисов текущей фазы
# 3. Запуск сервисов фазы
# 4. Gateway
```

## Переменные окружения (шаблон)

```env
APP_ENV=local
LOG_LEVEL=debug

POSTGRES_URL=postgres://socset:socset@localhost:5432/<service_db>?sslmode=disable
KAFKA_BROKERS=localhost:9092
REDIS_URL=redis://localhost:6379/0
ELASTICSEARCH_URL=http://localhost:9200
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minio
MINIO_SECRET_KEY=minio123
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
JWT_ISSUER=socset-auth
```

Файл `.env.example` будет добавлен с кодом; `.env` в git не коммитится.

## Миграции

- Инструмент: `golang-migrate` или `goose`
- Путь: `services/<name>/migrations/`
- Правило: только вперёд на shared dev; down — для локальных экспериментов

## Проверка здоровья

Каждый сервис:

- `GET /healthz` — liveness
- `GET /readyz` — готовность (DB/Kafka/ES)

## Тесты

| Уровень | Что |
|---------|-----|
| Unit | доменная логика |
| Integration | testcontainers: Postgres/Kafka/ES |
| E2E | по фазам через gateway |

На фазе 1 достаточно unit + integration auth/profile.

## Полезные локальные URL (ориентир)

| Сервис | URL |
|--------|-----|
| Gateway | http://localhost:8080 |
| Jaeger UI | http://localhost:16686 |
| MinIO Console | http://localhost:9001 |
| Kibana/ES | http://localhost:9200 |
| Kafka UI | http://localhost:8085 (если добавим) |

## Соглашения по коду Go

- `internal/` для приватного кода сервиса
- `cmd/<service>/main.go` — точка входа
- Конфиг через env
- Контекст со `request_id` / trace на все outbound вызовы
- Не логировать пароли и токены
