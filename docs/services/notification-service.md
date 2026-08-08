# notification-service

## Назначение

In-app уведомления (и задел под push) на основе доменных событий Kafka.

## Ответственность

- Consume доменных событий
- Создание уведомлений для пользователей
- API списка / mark read / unread count
- (Позже) интеграция WebPush / FCM

## Не делает

- Бизнес-логику дружбы/постов
- Гарантию доставки real-time чата (чат — отдельный канал)

## Модель данных (Postgres)

```text
notifications
  id            UUID PK
  user_id       UUID NOT NULL   -- получатель
  type          TEXT NOT NULL   -- friend_requested|post_liked|comment_created|message_sent|...
  actor_id      UUID            -- кто инициировал
  entity_type   TEXT            -- post|comment|user|message|...
  entity_id     UUID
  payload       JSONB
  read_at       TIMESTAMPTZ NULL
  created_at    TIMESTAMPTZ

processed_events
  event_id      UUID PK         -- идемпотентность
  processed_at  TIMESTAMPTZ
```

## Публичное API

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/v1/notifications` | Список (cursor) |
| GET | `/v1/notifications/unread-count` | Счётчик |
| POST | `/v1/notifications/{id}/read` | Прочитать одно |
| POST | `/v1/notifications/read-all` | Прочитать все |

## Consumed events

| Событие | Уведомление |
|---------|-------------|
| `friend.requested` | «X хочет добавить вас в друзья» → `to_user_id` |
| `friend.accepted` | «X принял заявку» → `from_user_id` |
| `post.liked` | «X лайкнул ваш пост» → author поста (не себе) |
| `comment.created` | «X прокомментировал» → author поста (не себе) |
| `message.sent` | «Новое сообщение» → peer (если не отправитель) |

Не создавать уведомление, если `actor_id == recipient_id`.

## Real-time (опционально)

- Отдельный WS `/ws/notifications` или SSE
- Либо клиент поллит unread-count на MVP

## Зависимости

- Postgres
- Kafka (consumer group `notification-service`)

## Фаза появления

Фаза 3.
