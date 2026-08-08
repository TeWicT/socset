# search-service

## Назначение

Поиск людей и постов. Elasticsearch — read-модель, не source of truth.

## Ответственность

- Индексация документов из Kafka
- HTTP API поиска
- Переиндексация (admin/job) при необходимости

## Не делает

- CRUD профилей/постов
- Авторизацию видимости сложных ACL на MVP (базово: публичные профили/посты; private — фильтр по флагу)

## Индексы Elasticsearch

### `users_v1`

```json
{
  "mappings": {
    "properties": {
      "user_id": { "type": "keyword" },
      "username": { "type": "text", "fields": { "keyword": { "type": "keyword" } } },
      "display_name": { "type": "text" },
      "is_private": { "type": "boolean" },
      "updated_at": { "type": "date" }
    }
  }
}
```

### `posts_v1`

```json
{
  "mappings": {
    "properties": {
      "post_id": { "type": "keyword" },
      "author_id": { "type": "keyword" },
      "wall_user_id": { "type": "keyword" },
      "body": { "type": "text" },
      "created_at": { "type": "date" }
    }
  }
}
```

Алиасы: `users`, `posts` → текущая версия индекса.

## Consumed events

| Событие | Действие |
|---------|----------|
| `profile.created` / `profile.updated` | upsert user doc |
| `user.registered` | опционально stub до profile |
| `post.created` | index post |
| `post.deleted` | delete post doc |
| `user.deleted` | delete user doc (+ посты — политика) |

## Публичное API

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/v1/search/users?q=&limit=` | Поиск людей |
| GET | `/v1/search/posts?q=&limit=` | Поиск постов |

Ответ — документы из ES + `user_id`/`post_id` для дальнейшей догрузки.

## Идемпотентность

- `event_id` в `processed_events` (Postgres или отдельный store)
- upsert по `user_id` / `post_id`

## Зависимости

- Elasticsearch
- Kafka
- Postgres (опционально для outbox offset / processed_events)

## Фаза появления

Фаза 3.
