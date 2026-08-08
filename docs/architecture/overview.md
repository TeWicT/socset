# Архитектура

## Цель

Мини-соцсеть с микросервисной архитектурой, приближенной к прод-практикам бигтеха:

- один bounded context — один сервис — своя БД;
- синхронный edge (HTTP через gateway);
- асинхронная связность через Kafka;
- read-модели для поиска (Elasticsearch);
- сквозная наблюдаемость (OpenTelemetry).

## Контекст (C4 L1)

```text
                   ┌─────────────┐
                   │   Clients   │
                   │ Web / Mobile│
                   └──────┬──────┘
                          │ HTTPS / WSS
                          ▼
                   ┌─────────────┐
                   │ api-gateway │
                   └──────┬──────┘
          ┌───────────────┼───────────────┐
          │               │               │
          ▼               ▼               ▼
     auth-service   profile / social   chat-service
                    graph / post       (WS upgrade)
                          │
                          ▼
              notification / search
                   (Kafka consumers)
```

## Контейнеры (C4 L2)

| Контейнер | Технология | Роль |
|-----------|------------|------|
| api-gateway | Go | Единая точка входа HTTP; JWT; rate limit; `X-Request-Id` |
| auth-service | Go + Postgres | Identity, credentials, выдача токенов |
| profile-service | Go + Postgres + MinIO | Публичный профиль, настройки, аватар |
| social-graph-service | Go + Postgres | Заявки и дружба |
| post-service | Go + Postgres | Посты стены, лайки, комментарии |
| chat-service | Go + Postgres + Redis | Диалоги, сообщения, WS presence |
| notification-service | Go + Postgres | In-app уведомления |
| search-service | Go + Elasticsearch | Поиск людей и постов |
| Kafka | Apache Kafka | Доменные события |
| Redis | Redis | Кэш / rate limit / presence |
| MinIO | S3 API | Медиа (аватары) |
| OTel Collector | OpenTelemetry | Трейсы, метрики, логи |

## Принципы

1. **Database per service** — нет shared schema между сервисами.
2. **Authz на gateway + сервисе** — gateway проверяет JWT; сервис доверяет `user_id` из доверенных заголовков/claims только от gateway (mTLS или internal network).
3. **Синхронно только там, где нужен ответ здесь и сейчас** — CRUD через HTTP/gRPC.
4. **Асинхронно для fan-out** — уведомления, индексация поиска, eventual consistency счётчиков.
5. **Transactional outbox** — публикация в Kafka после коммита в БД источника.
6. **Идемпотентные consumer'ы** — повтор сообщения не ломает состояние.
7. **ES — не source of truth** — только проекция для поиска.

## Основные пользовательские потоки

### Регистрация

```text
Client → Gateway → auth-service (create user + credentials)
                 → Kafka: user.registered
                 → profile-service создаёт пустой профиль
                 → search-service индексирует пользователя (после появления профиля)
```

### Просмотр профиля + стена

```text
Client → Gateway → profile-service (карточка)
                 → post-service (посты, лайки, комментарии)
                 → social-graph-service (статус дружбы относительно viewer)
```

Агрегацию UI может делать клиент или BFF-слой gateway (composition). На MVP — параллельные запросы с клиента через gateway.

### Заявка в друзья

```text
Client → Gateway → social-graph-service
                 → Kafka: friend.requested
                 → notification-service
```

### Поиск

```text
Client → Gateway → search-service → Elasticsearch
```

### Чат

```text
Client ⇄ WSS → Gateway (или прямой WS endpoint) ⇄ chat-service
chat-service → Kafka: message.sent → notification-service
```

## Синхронные vs асинхронные зависимости

| От | К | Тип | Зачем |
|----|---|-----|-------|
| gateway | все HTTP-сервисы | sync | API |
| auth | — | publishes `user.*` | создание identity |
| profile | MinIO | sync | аватары |
| social-graph | — | publishes `friend.*` | уведомления |
| post | — | publishes `post.*` / `comment.*` | search + notifications |
| chat | Redis | sync | presence / pubsub внутри инстансов |
| chat | — | publishes `message.sent` | notifications |
| notification | Kafka | async consume | доставка |
| search | Kafka | async consume | индексация |
| search | ES | sync write/read | индекс / запрос |

**Запрещено на старте:** sync-вызовы post → profile «на каждый лайк», search → post как source of truth, shared DB.

## Идентичность пользователя

- Канонический идентификатор: `user_id` (UUID), выдаётся **auth-service** при регистрации.
- Остальные сервисы хранят только `user_id`, не пароли и не PII credentials.
- JWT access token содержит `sub = user_id`, `exp`, опционально `sid` (session).

## Ошибки и согласованность

- Клиент всегда получает синхронный ответ от сервиса-владельца данных.
- Проекции (search, notifications) обновляются eventually.
- При отставании consumer'а UI может кратковременно показывать устаревший поиск/бейдж уведомлений — это ожидаемо.

## Дальше

- [Владение данными](data-ownership.md)
- [ADR](decisions.md)
- [Каталог событий](../events/catalog.md)
