# Инфраструктурный стек

## Компоненты

| Компонент | Назначение | Кто использует |
|-----------|------------|----------------|
| PostgreSQL | Source of truth сервисов | auth, profile, social-graph, post, chat, notification (+ outbox) |
| Apache Kafka (+ ZooKeeper/KRaft) | Доменные события | все producers/consumers |
| Elasticsearch | Поисковые индексы | search-service |
| Redis | Rate limit, presence, pub/sub chat | gateway, chat |
| MinIO | S3-совместимое хранилище медиа | profile-service |
| OpenTelemetry Collector | Приём телеметрии | все сервисы |
| Jaeger / Grafana Tempo | Трейсы | ops |
| Prometheus + Grafana | Метрики | ops |
| (опц.) Redpanda Console / AKHQ | UI для Kafka | dev |

На локалке допустим упрощённый набор: Postgres, Kafka, ES, Redis, MinIO, Jaeger all-in-one.

## Сеть сервисов (логически)

```text
                 Internet
                     │
               api-gateway :8080
                     │
     ┌───────────────┼────────────────┐
     │ internal docker network        │
     ▼               ▼                ▼
  services...     Kafka           Redis/ES/MinIO/PG
```

Снаружи открыты: gateway `:8080`, (опц.) MinIO console. Остальное — internal.

## Postgres: изоляция

Варианты:

1. **Один инстанс Postgres, разные databases** — удобно для локалки (`auth`, `profile`, …).
2. **Отдельные инстансы** — ближе к проду, тяжелее локально.

Решение для MVP локалки: один Postgres, database per service. В ADR это всё ещё database-per-service.

## Kafka

- Партиции: старт с 3 на топик (локально можно 1).
- Retention: 7 дней на dev.
- Consumer groups: `profile-service`, `search-service`, `notification-service`, …

## Elasticsearch

- Один узел на локалке.
- Индексы с алиасами `users`, `posts`.
- Память: ограничить heap в compose (например 512MB–1GB).

## Redis

- DB0: gateway rate limit
- DB1: chat presence / pubsub (или key prefixes в одном DB)

## MinIO

- Bucket: `avatars`
- Prefixed keys: `{user_id}/{object_id}`
- Публичный read через gateway/presigned URL

## Observability

Минимум на сервис:

- structured JSON logs: `level`, `msg`, `request_id`, `trace_id`, `user_id`
- OTel HTTP/gRPC instrumentation
- метрики RED: rate, errors, duration

## Секреты

- `.env` локально, не коммитить
- Позже: Docker secrets / SOPS / vault-like — не блокирует MVP

## Целевой `deploy/docker-compose.yml` (ориентир)

Сервисы приложения + зависимости:

- postgres, kafka, redis, minio, elasticsearch, jaeger, otel-collector
- api-gateway, auth, profile, … (подключаются по фазам)

Файл появится в репозитории на фазе 1 реализации.
