# api-gateway

## Назначение

Единая точка входа для HTTP(S) клиентов. Не владеет бизнес-данными.

## Ответственность

- Маршрутизация на downstream-сервисы
- Проверка JWT access token (authn)
- Прокидывание `X-Request-Id` / `X-User-Id` во внутренние запросы
- Rate limiting (по IP и/или `user_id`, Redis)
- CORS, базовые security headers
- (Опционально позже) BFF composition

## Не делает

- Авторизацию бизнес-правил (это сервисы)
- Хранение пользователей/постов
- Долгую бизнес-логику

## API (edge)

Публичные маршруты проксируются, например:

| Метод | Путь | Upstream |
|-------|------|----------|
| POST | `/api/v1/auth/register` | auth-service |
| POST | `/api/v1/auth/login` | auth-service |
| POST | `/api/v1/auth/refresh` | auth-service |
| GET/PATCH | `/api/v1/profiles/...` | profile-service |
| * | `/api/v1/friends/...` | social-graph-service |
| * | `/api/v1/posts/...` | post-service |
| GET | `/api/v1/search/...` | search-service |
| * | `/api/v1/notifications/...` | notification-service |
| GET/POST | `/api/v1/chats/...` | chat-service |
| WS | `/ws/chats` | chat-service |

Точный контракт — в [docs/api](../api/README.md).

## Authn

1. Извлекает `Authorization: Bearer <access_jwt>`.
2. Валидирует подпись и `exp` (JWKS или shared secret от auth).
3. При успехе добавляет внутренний заголовок `X-User-Id: <sub>`.
4. Публичные эндпоинты (register/login, публичный профиль) — без JWT.

## Конфиг (ориентир)

```yaml
listen_addr: ":8080"
jwt:
  jwks_url: "http://auth-service:8081/.well-known/jwks.json"
redis_url: "redis://redis:6379/0"
rate_limit:
  rps: 50
  burst: 100
upstreams:
  auth: "http://auth-service:8081"
  profile: "http://profile-service:8082"
  # ...
```

## Зависимости

- Redis — rate limit
- auth-service — JWKS / ключи
- все HTTP-сервисы — upstream

## Наблюдаемость

- Прокидывает и генерирует `X-Request-Id`
- Спан на каждый входящий запрос + child span на upstream
- Метрики: RPS, latency, 4xx/5xx, rate-limit hits

## Фаза появления

Фаза 1 (см. [roadmap](../roadmap.md)).
