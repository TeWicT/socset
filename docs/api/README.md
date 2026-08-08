# API

Публичный контракт клиентов. Источник правды по мере появления кода — OpenAPI в `schemas/openapi/`. Здесь — сводный черновик для проектирования.

## Общие правила

| Правило | Значение |
|---------|----------|
| Base URL | `https://<host>/api` (локально `http://localhost:8080/api`) |
| Версия | `/v1` |
| Формат | JSON (`Content-Type: application/json`) |
| Auth | `Authorization: Bearer <access_token>` |
| Request id | Клиент может прислать `X-Request-Id`; иначе gateway генерирует |
| Ошибки | Единый envelope |
| Пагинация | Cursor-based: `?cursor=&limit=` |

### Error envelope

```json
{
  "error": {
    "code": "friend_request_already_exists",
    "message": "Human readable message",
    "details": {}
  }
}
```

HTTP-статусы: `400` validation, `401` unauthenticated, `403` forbidden, `404` not found, `409` conflict, `429` rate limit, `500` internal.

### Cursor page

```json
{
  "items": [],
  "next_cursor": "opaque-or-null"
}
```

## Auth — `/v1/auth`

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| POST | `/register` | no | Регистрация |
| POST | `/login` | no | Логин |
| POST | `/refresh` | no | Обновление токенов |
| POST | `/logout` | refresh/access | Revoke refresh |

## Profiles — `/v1/profiles`

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| GET | `/{user_id}` | optional | Публичный профиль |
| GET | `/me` | yes | Свой профиль |
| PATCH | `/me` | yes | Обновить |
| POST | `/me/avatar` | yes | Загрузить аватар |

## Friends — `/v1/friends`

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| GET | `/` | yes | Список друзей |
| GET | `/status/{user_id}` | yes | Статус отношений |
| POST | `/requests` | yes | Создать заявку |
| GET | `/requests/incoming` | yes | Входящие |
| GET | `/requests/outgoing` | yes | Исходящие |
| POST | `/requests/{id}/accept` | yes | Принять |
| POST | `/requests/{id}/reject` | yes | Отклонить |
| DELETE | `/requests/{id}` | yes | Отменить |
| DELETE | `/{user_id}` | yes | Удалить из друзей |

## Posts / Wall — `/v1`

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| GET | `/walls/{user_id}/posts` | optional* | Стена |
| POST | `/walls/{user_id}/posts` | yes | Создать пост |
| DELETE | `/posts/{id}` | yes | Удалить |
| POST | `/posts/{id}/likes` | yes | Лайк |
| DELETE | `/posts/{id}/likes` | yes | Анлайк |
| GET | `/posts/{id}/comments` | optional* | Комментарии |
| POST | `/posts/{id}/comments` | yes | Коммент |
| DELETE | `/comments/{id}` | yes | Удалить коммент |

\*с учётом `is_private` и дружбы — уточняется в фазе 2.

## Search — `/v1/search`

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| GET | `/users?q=` | yes | Поиск людей |
| GET | `/posts?q=` | yes | Поиск постов |

## Notifications — `/v1/notifications`

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| GET | `/` | yes | Список |
| GET | `/unread-count` | yes | Счётчик |
| POST | `/{id}/read` | yes | Прочитать |
| POST | `/read-all` | yes | Прочитать все |

## Chats — `/v1/chats` + WS

| Метод | Путь | Auth | Описание |
|-------|------|------|----------|
| GET | `/` | yes | Диалоги |
| POST | `/` | yes | Создать/получить DM |
| GET | `/{id}/messages` | yes | История |
| WS | `/ws/chats` | yes | Real-time |

Детали WS — в [chat-service](../services/chat-service.md).

## Связанные документы

- [Сервисы](../services/)
- [События](../events/catalog.md)
- OpenAPI файлы — появятся в `schemas/openapi/` на фазах реализации
