# api-gateway

## Назначение

Единая точка входа для HTTP(S) клиентов. Не владеет бизнес-данными.

## Ответственность

- Маршрутизация на downstream-сервисы
- Проверка JWT access token (authn), в т.ч. Redis denylist по `jti`
- Прокидывание `X-Request-Id` / `X-User-Id` во внутренние запросы
- Rate limiting (по IP и/или `user_id`, Redis) — позже
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
| POST | `/api/v1/auth/register` | auth-service (gRPC) |
| POST | `/api/v1/auth/login` | auth-service (gRPC) |
| POST | `/api/v1/auth/refresh` | auth-service (gRPC) |
| POST | `/api/v1/auth/logout` | auth-service (gRPC) + Redis denylist; **требует JWT** |
| GET | `/api/v1/auth/me` | gateway (JWT → `sub`) |
| GET | `/api/v1/profiles/me` | profile-service (gRPC), JWT |
| PATCH | `/api/v1/profiles/me` | profile-service (gRPC), JWT |
| GET | `/api/v1/profiles/{user_id}` | profile-service (gRPC) |
| * | `/api/v1/friends/...` | social-graph-service |
| * | `/api/v1/posts/...` | post-service |
| GET | `/api/v1/search/...` | search-service |
| * | `/api/v1/notifications/...` | notification-service |
| GET/POST | `/api/v1/chats/...` | chat-service |
| WS | `/ws/chats` | chat-service |

Точный контракт — в [docs/api](../api/README.md).

## Authn

1. Извлекает `Authorization: Bearer <access_jwt>`.
2. Валидирует подпись, `exp`, наличие на `jti` (shared secret HS256 на MVP).
3. Проверяет Redis denylist: ключ `access:deny:{jti}`. Если есть — `401` (токен отозван при logout).
4. При успехе кладёт в request context `user_id` (`sub`), `jti`, `exp` (для logout).
5. Публичные эндпоинты (register/login/refresh, публичный профиль) — без JWT.

### Logout и denylist

1. `POST /api/v1/auth/logout` защищён JWT; в body — `refresh_token`.
2. Gateway вызывает auth-service: revoke refresh-сессии.
3. Затем `SET access:deny:{jti} revoked` в Redis с TTL = время до `exp` access.
4. Пока TTL не истёк, тот же access на защищённых роутах получает `401`.

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

- Redis — denylist access (`jti`); позже rate limit
- auth-service — gRPC, shared `JWT_SECRET`
- profile-service — gRPC
- остальные сервисы — по мере фаз

## Наблюдаемость

- Прокидывает и генерирует `X-Request-Id`
- Спан на каждый входящий запрос + child span на upstream
- Метрики: RPS, latency, 4xx/5xx, rate-limit hits

## Фаза появления

Фаза 1 (см. [roadmap](../roadmap.md)).
