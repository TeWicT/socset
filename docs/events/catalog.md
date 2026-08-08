# Каталог событий Kafka

Единый источник правды по доменным событиям. Перед добавлением нового события — обновить этот файл.

## Соглашения

| Правило | Описание |
|---------|----------|
| Имя топика | `socset.<domain>.<entity>.<event>` или короче `socset.<domain>.events` с полем `type` |
| На MVP | Один топик на домен: `socset.users.events`, `socset.profiles.events`, `socset.friends.events`, `socset.posts.events`, `socset.chats.events` |
| Ключ сообщения | `user_id` / `post_id` / `conversation_id` — для упорядочивания в партиции |
| Формат | JSON (Avro/Protobuf — позже) |
| Обязательные поля envelope | `event_id`, `event_type`, `event_version`, `occurred_at`, `producer` |
| Доставка | At-least-once → consumer'ы идемпотентны |
| Публикация | Transactional outbox в сервисе-производителе |

### Envelope

```json
{
  "event_id": "018f...",
  "event_type": "profile.updated",
  "event_version": 1,
  "occurred_at": "2026-08-08T12:00:00Z",
  "producer": "profile-service",
  "payload": {}
}
```

## Топики

| Топик | Производитель | Consumer'ы |
|-------|---------------|------------|
| `socset.users.events` | auth-service | profile-service, search-service |
| `socset.profiles.events` | profile-service | search-service |
| `socset.friends.events` | social-graph-service | notification-service |
| `socset.posts.events` | post-service | notification-service, search-service |
| `socset.chats.events` | chat-service | notification-service |

## События

### users

#### `user.registered` v1

```json
{
  "user_id": "uuid",
  "username": "alice",
  "email": "alice@example.com"
}
```

#### `user.deleted` v1 (позже)

```json
{ "user_id": "uuid" }
```

### profiles

#### `profile.created` v1 / `profile.updated` v1

```json
{
  "user_id": "uuid",
  "username": "alice",
  "display_name": "Alice",
  "is_private": false,
  "avatar_url": "https://..."
}
```

`username` дублируется в событии для удобства индексации (денормализация в событии допустима).

### friends

#### `friend.requested` v1

```json
{
  "request_id": "uuid",
  "from_user_id": "uuid",
  "to_user_id": "uuid"
}
```

#### `friend.accepted` v1

```json
{
  "request_id": "uuid",
  "from_user_id": "uuid",
  "to_user_id": "uuid"
}
```

#### `friend.rejected` v1 / `friend.canceled` v1 / `friend.removed` v1

```json
{
  "from_user_id": "uuid",
  "to_user_id": "uuid"
}
```

### posts

#### `post.created` v1

```json
{
  "post_id": "uuid",
  "author_id": "uuid",
  "wall_user_id": "uuid",
  "body": "text",
  "created_at": "RFC3339"
}
```

#### `post.deleted` v1

```json
{ "post_id": "uuid", "author_id": "uuid" }
```

#### `post.liked` v1 / `post.unliked` v1

```json
{
  "post_id": "uuid",
  "post_author_id": "uuid",
  "liker_id": "uuid"
}
```

#### `comment.created` v1

```json
{
  "comment_id": "uuid",
  "post_id": "uuid",
  "post_author_id": "uuid",
  "author_id": "uuid",
  "body": "text"
}
```

#### `comment.deleted` v1

```json
{ "comment_id": "uuid", "post_id": "uuid" }
```

### chats

#### `message.sent` v1

```json
{
  "message_id": "uuid",
  "conversation_id": "uuid",
  "sender_id": "uuid",
  "recipient_ids": ["uuid"],
  "body_preview": "truncated..."
}
```

Полный текст в Kafka на MVP допустим; в более строгой модели — только preview + id.

## Outbox (обязательный паттерн)

В БД производителя:

```text
outbox_events
  id            UUID PK   -- = event_id
  topic         TEXT
  partition_key TEXT
  payload       JSONB
  created_at    TIMESTAMPTZ
  published_at  TIMESTAMPTZ NULL
```

Relay (отдельная горутина/процесс) читает unpublished и пишет в Kafka.

## Версионирование

- Additive changes → тот же `event_version`, опциональные поля.
- Breaking → `event_version + 1` или новый `event_type`; consumer'ы поддерживают N и N-1 на период миграции.

## Что не кладём в Kafka

- Сырые пароли, refresh tokens
- Полные бинарные файлы (аватары — только URL/key)
