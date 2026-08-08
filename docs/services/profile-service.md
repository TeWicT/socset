# profile-service

## Назначение

Публичный и приватный профиль пользователя, настройки, аватар.

## Ответственность

- CRUD профиля
- Загрузка/смена аватара (MinIO)
- Настройки приватности (базовые)
- Batch-чтение профилей для обогащения ленты
- События `profile.created` / `profile.updated`

## Не делает

- Пароли и токены
- Дружбу и заявки
- Посты

## Модель данных (Postgres)

```text
profiles
  user_id       UUID PK          -- = auth user_id
  display_name  TEXT NOT NULL
  bio           TEXT
  avatar_key    TEXT             -- object key в MinIO
  is_private    BOOLEAN DEFAULT false
  created_at    TIMESTAMPTZ
  updated_at    TIMESTAMPTZ
```

Создание профиля: consumer `user.registered` **или** sync-вызов из auth на MVP. Предпочтительно: Kafka `user.registered` → profile создаёт запись с `display_name = username`.

## Публичное API

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/v1/profiles/{user_id}` | Публичный профиль |
| GET | `/v1/profiles/me` | Свой профиль |
| PATCH | `/v1/profiles/me` | Обновить bio / display_name / privacy |
| POST | `/v1/profiles/me/avatar` | Upload avatar (multipart) |
| GET | `/internal/v1/profiles:batch` | `?ids=` batch (internal) |

## События

| Событие | Когда |
|---------|-------|
| `profile.created` | Профиль создан |
| `profile.updated` | Изменены индексируемые поля |

Consumer'ы: search-service (индекс людей), опционально кэши.

## Зависимости

- Postgres
- MinIO
- Kafka (consume `user.registered`, produce `profile.*`)

## Фаза появления

Фаза 1 (базовый профиль) → Фаза 2 (аватар, privacy).
