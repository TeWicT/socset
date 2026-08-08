# social-graph-service

## Назначение

Социальный граф: заявки в друзья и принятая дружба.

## Ответственность

- Отправить / отменить / принять / отклонить заявку
- Список друзей
- Статус отношений относительно viewer (`none` | `outgoing_request` | `incoming_request` | `friends`)
- События `friend.requested` / `friend.accepted` / `friend.rejected` / `friend.removed`

## Не делает

- Профили и аватары (только `user_id`)
- Ленту постов

## Модель данных (Postgres)

```text
friend_requests
  id            UUID PK
  from_user_id  UUID NOT NULL
  to_user_id    UUID NOT NULL
  status        TEXT NOT NULL  -- pending|accepted|rejected|canceled
  created_at    TIMESTAMPTZ
  updated_at    TIMESTAMPTZ
  UNIQUE (from_user_id, to_user_id)

friendships
  user_id_a     UUID NOT NULL
  user_id_b     UUID NOT NULL
  created_at    TIMESTAMPTZ
  PRIMARY KEY (user_id_a, user_id_b)
  CHECK (user_id_a < user_id_b)   -- каноническая пара
```

## Публичное API

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/v1/friends/requests` | Создать заявку `{ to_user_id }` |
| POST | `/v1/friends/requests/{id}/accept` | Принять |
| POST | `/v1/friends/requests/{id}/reject` | Отклонить |
| DELETE | `/v1/friends/requests/{id}` | Отменить свою pending |
| GET | `/v1/friends` | Список друзей текущего пользователя |
| GET | `/v1/friends/requests/incoming` | Входящие |
| GET | `/v1/friends/requests/outgoing` | Исходящие |
| GET | `/v1/friends/status/{user_id}` | Статус относительно me |
| DELETE | `/v1/friends/{user_id}` | Удалить из друзей |

## Инварианты

- Нельзя добавить себя
- Одна активная pending-заявка на пару
- Accept идемпотентен
- Дружба симметрична (одна запись-пара)

## События

См. [catalog](../events/catalog.md): `friend.*` → notification-service.

## Зависимости

- Postgres
- Kafka (outbox)
- Нет обязательного sync на profile (клиент сам подтянет имена)

## Фаза появления

Фаза 2.
