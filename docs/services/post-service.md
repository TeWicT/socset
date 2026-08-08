# post-service

## Назначение

Стена профиля: посты, лайки, комментарии.

## Ответственность

- Создать / удалить / получить посты пользователя (стена)
- Лайк / анлайк
- Комментарии к посту
- События `post.created`, `post.deleted`, `post.liked`, `post.unliked`, `comment.created`, `comment.deleted`

## Не делает

- Хранение полного профиля автора
- Ленту «новости друзей» как отдельный продукт на MVP (можно добавить позже как read-model)

## Модель данных (Postgres)

```text
posts
  id            UUID PK
  author_id     UUID NOT NULL
  wall_user_id  UUID NOT NULL   -- на чьей стене
  body          TEXT NOT NULL
  created_at    TIMESTAMPTZ
  deleted_at    TIMESTAMPTZ NULL

likes
  post_id       UUID NOT NULL
  user_id       UUID NOT NULL
  created_at    TIMESTAMPTZ
  PRIMARY KEY (post_id, user_id)

comments
  id            UUID PK
  post_id       UUID NOT NULL
  author_id     UUID NOT NULL
  body          TEXT NOT NULL
  created_at    TIMESTAMPTZ
  deleted_at    TIMESTAMPTZ NULL

-- опционально денормализация
post_counters
  post_id       UUID PK
  likes_count   BIGINT
  comments_count BIGINT
```

Счётчики можно обновлять транзакционно в том же сервисе или eventually через внутренние события — на MVP транзакционно вместе с like/comment.

## Публичное API

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/v1/walls/{user_id}/posts` | Посты стены (cursor pagination) |
| POST | `/v1/walls/{user_id}/posts` | Создать пост (на своей стене или по правилам) |
| DELETE | `/v1/posts/{id}` | Удалить свой пост |
| POST | `/v1/posts/{id}/likes` | Лайк |
| DELETE | `/v1/posts/{id}/likes` | Анлайк |
| GET | `/v1/posts/{id}/comments` | Комментарии |
| POST | `/v1/posts/{id}/comments` | Добавить комментарий |
| DELETE | `/v1/comments/{id}` | Удалить свой комментарий |

### Правила стены (MVP)

- Автор может писать на свою стену.
- Друг может писать на стену друга (проверка через sync-вызов social-graph **или** временно «только своя стена» до появления graph). Предпочтительно: на фазе 2 — только своя стена; на фазе 2.1 — проверка friendship.

## События

| Событие | Consumer'ы |
|---------|------------|
| `post.created` / `post.deleted` | search, (опц.) notification |
| `post.liked` | notification (владельцу поста) |
| `comment.created` | notification, (опц.) search |

## Зависимости

- Postgres
- Kafka
- Опционально sync: social-graph `GET /friends/status/{id}` для ACL стены

## Фаза появления

Фаза 2.
