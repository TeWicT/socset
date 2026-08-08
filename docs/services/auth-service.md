# auth-service

## Назначение

Identity-провайдер: регистрация, аутентификация, выдача и ротация токенов.

## Ответственность

- Регистрация пользователя (`user_id`, email/username, password hash)
- Login / logout
- Access JWT + refresh token
- Revoke refresh / logout all sessions
- Публикация JWKS (если asymmetric keys)
- Событие `user.registered` (и при необходимости `user.deleted`)

## Не делает

- Профиль (имя, bio, аватар) — это profile-service
- Друзей, посты, чаты

## Модель данных (Postgres)

```text
users
  id            UUID PK
  email         CITEXT UNIQUE NOT NULL
  username      CITEXT UNIQUE NOT NULL
  password_hash TEXT NOT NULL
  created_at    TIMESTAMPTZ
  updated_at    TIMESTAMPTZ

refresh_sessions
  id            UUID PK
  user_id       UUID FK
  token_hash    TEXT NOT NULL
  expires_at    TIMESTAMPTZ
  revoked_at    TIMESTAMPTZ NULL
  user_agent    TEXT
  created_at    TIMESTAMPTZ
```

## Публичное API

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/v1/auth/register` | Регистрация → tokens |
| POST | `/v1/auth/login` | Логин → tokens |
| POST | `/v1/auth/refresh` | Обмен refresh → новая пара |
| POST | `/v1/auth/logout` | Revoke refresh |
| GET | `/.well-known/jwks.json` | Ключи проверки JWT |

### Register (черновик)

```json
{
  "email": "user@example.com",
  "username": "alice",
  "password": "********"
}
```

Ответ: `access_token`, `refresh_token`, `expires_in`, `user_id`.

## События

| Событие | Когда |
|---------|-------|
| `user.registered` | Успешная регистрация |
| `user.deleted` | (позже) удаление аккаунта |

Payload `user.registered`:

```json
{
  "event_id": "uuid",
  "occurred_at": "RFC3339",
  "user_id": "uuid",
  "username": "alice",
  "email": "user@example.com"
}
```

Email в событии — осознанный компромисс для MVP; в более строгой модели PII минимизируют, а profile создаётся sync-вызовом/отдельным полем без email.

## Безопасность

- Argon2id или bcrypt для паролей
- Refresh — только hash в БД
- Короткий TTL access (например 15m), refresh — дни
- Rate limit brute-force на login (gateway + локально)

## Зависимости

- Postgres
- Kafka (outbox)
- Нет sync-зависимостей от других бизнес-сервисов

## Фаза появления

Фаза 1.
