# Владение данными

Каждый сервис — единственный владелец своих сущностей. Чужие данные читаются только через API владельца или через собственную проекцию, построенную из Kafka-событий.

## Матрица владения

| Сущность | Владелец | Хранилище | Кто ещё использует |
|----------|----------|-----------|--------------------|
| credentials, refresh sessions | auth-service | Postgres | — |
| user_id (identity) | auth-service | Postgres | все (как FK-ссылка) |
| profile (display name, bio, avatar_url, settings) | profile-service | Postgres + MinIO | search (проекция), клиент |
| friendship, friend_request | social-graph-service | Postgres | notification (события) |
| post, like, comment | post-service | Postgres | search (проекция), notification |
| conversation, message | chat-service | Postgres | notification |
| presence / online | chat-service | Redis | — |
| notification | notification-service | Postgres | — |
| search documents | search-service | Elasticsearch | — (read model) |

## Что нельзя делать

- Читать/писать чужую БД напрямую.
- Дублировать пароли или refresh-токены вне auth.
- Считать Elasticsearch источником правды для профилей/постов.
- Хранить полный профиль внутри post/chat «на всякий случай» — только `user_id` (+ опциональный короткий кэш с TTL, инвалидируемый по `user.updated` / `profile.updated`).

## Ссылки между сервисами

Все межсервисные ссылки — по `user_id` / `post_id` / `conversation_id` (UUID).

Пример: комментарий в post-service:

```text
comment {
  id,
  post_id,      // владеет post-service
  author_id,    // ссылка на auth/profile user_id
  body,
  created_at
}
```

Для отображения имени автора клиент (или composition на gateway) запрашивает profile-service batch-эндпоинтом `GET /profiles?ids=...`.

## Batch и кэш

- Profile: `GET /internal/profiles:batch` для обогащения ленты/комментов.
- Краткий кэш display name/avatar в Redis на gateway или в post-service — опционально, TTL 1–5 мин, сброс по событию `profile.updated`.

## Миграции

- Миграции только внутри сервиса-владельца.
- Версионирование схем событий — в [каталоге событий](../events/catalog.md); breaking changes только через новый major или новый топик.
