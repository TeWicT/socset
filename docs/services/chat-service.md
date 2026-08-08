# chat-service

## Назначение

Диалоги и сообщения в реальном времени по WebSocket.

## Ответственность

- Создание/получение диалогов (1:1 на MVP; группы — позже)
- История сообщений
- WebSocket: отправка/получение сообщений online
- Presence (online/offline) в Redis
- Pub/Sub между инстансами chat через Redis
- Событие `message.sent` для notification-service

## Не делает

- Push-доставку (это notification)
- Поиск по сообщениям на MVP

## Модель данных (Postgres)

```text
conversations
  id            UUID PK
  type          TEXT NOT NULL  -- dm
  created_at    TIMESTAMPTZ

conversation_members
  conversation_id UUID NOT NULL
  user_id         UUID NOT NULL
  joined_at       TIMESTAMPTZ
  PRIMARY KEY (conversation_id, user_id)

messages
  id              UUID PK
  conversation_id UUID NOT NULL
  sender_id       UUID NOT NULL
  body            TEXT NOT NULL
  created_at      TIMESTAMPTZ
  deleted_at      TIMESTAMPTZ NULL
```

Для DM: уникальность пары участников (канонический ключ `min(user_a), max(user_b)`).

## HTTP API

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/v1/chats` | Список диалогов текущего пользователя |
| POST | `/v1/chats` | Создать/получить DM `{ peer_user_id }` |
| GET | `/v1/chats/{id}/messages` | История (cursor) |

## WebSocket

- Endpoint: `/ws/chats` (через gateway или напрямую)
- Auth: JWT при handshake (`Sec-WebSocket-Protocol` или query — предпочтительно header через gateway)
- Клиент подписывается на свои conversation ids автоматически по membership

### Клиент → сервер (черновик)

```json
{ "type": "message.send", "conversation_id": "...", "body": "hi", "client_msg_id": "..." }
{ "type": "typing.start", "conversation_id": "..." }
```

### Сервер → клиент

```json
{ "type": "message.created", "message": { "...": "..." } }
{ "type": "presence", "user_id": "...", "status": "online" }
```

`client_msg_id` — идемпотентность отправки.

## Масштабирование

```text
Client A ──WS──► chat-instance-1 ──Redis Pub/Sub──► chat-instance-2 ◄──WS── Client B
                      │
                      ▼
                   Postgres
```

## События

`message.sent` → notification (если peer offline или всегда для inbox badge).

## Зависимости

- Postgres
- Redis
- Kafka
- JWT validation (локально или через shared JWKS)

## Фаза появления

Фаза 4.
