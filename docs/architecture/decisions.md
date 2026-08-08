# Architecture Decision Records (ADR)

Формат: короткий ADR. Статусы: `Proposed` · `Accepted` · `Superseded` · `Deprecated`.

---

## ADR-001: Микросервисные границы

- **Статус:** Accepted
- **Контекст:** Нужен опыт, близкий к бигтеху: отдельные команды/репозитории логически, независимый деплой, своя БД.
- **Решение:** 8 сервисов — gateway, auth, profile, social-graph, post, chat, notification, search.
- **Последствия:** Больше операционной сложности; компенсируется поэтапным roadmap и docker-compose.

---

## ADR-002: Database per service

- **Статус:** Accepted
- **Контекст:** Shared DB превращает микросервисы в распределённый монолит.
- **Решение:** У каждого сервиса свой Postgres (отдельная БД/инстанс). ES только у search. Redis — инфраструктурный кэш, не SoT.
- **Последствия:** Нет кросс-сервисных JOIN; нужны API/события и eventual consistency.

---

## ADR-003: Kafka как шина доменных событий

- **Статус:** Accepted
- **Контекст:** Нужны уведомления, поиск, слабая связность.
- **Решение:** Доменные события публикуются в Kafka через transactional outbox. Consumer'ы идемпотентны.
- **Последствия:** Нужен outbox-relay; задержка проекций; обязательный event catalog.

---

## ADR-004: Elasticsearch только как read-модель

- **Статус:** Accepted
- **Контекст:** Хотим полноценный поиск как в проде.
- **Решение:** search-service слушает Kafka и индексирует документы. Source of truth остаётся в Postgres владельцев.
- **Последствия:** Возможна кратковременная рассинхронизация поиска.

---

## ADR-005: Внешний REST, внутренний gRPC (цель)

- **Статус:** Accepted
- **Контекст:** Бигтех часто разделяет public HTTP и internal RPC.
- **Решение:** Снаружи — REST/OpenAPI (+ WS для чата). Между сервисами цель — gRPC. На ранних фазах допустим internal HTTP, затем миграция на gRPC без смены границ.
- **Последствия:** Два стиля контрактов; явная версия internal API.

---

## ADR-006: JWT access + opaque refresh

- **Статус:** Accepted
- **Контекст:** Нужны stateless проверки на gateway и возможность ревокации сессий.
- **Решение:** Короткоживущий JWT access; refresh хранится в auth-service (hash), возможен logout/revoke.
- **Последствия:** Gateway валидирует JWT (JWKS/shared key); refresh только через auth.

---

## ADR-007: WebSocket для чата

- **Статус:** Accepted
- **Контекст:** Real-time сообщения.
- **Решение:** chat-service владеет WS-соединениями. Gateway либо проксирует WS со sticky, либо клиент подключается к отдельному WS endpoint chat-service за тем же JWT.
- **Последствия:** Нужна стратегия горизонтального скейла (Redis pub/sub между инстансами chat).

---

## ADR-008: Composition на клиенте (MVP)

- **Статус:** Accepted
- **Контекст:** Страница профиля требует данные из profile + post + social-graph.
- **Решение:** На MVP клиент делает параллельные запросы через gateway. BFF/aggregation в gateway — опционально позже.
- **Последствия:** Больше round-trips; проще сервисы.

---

## ADR-009: Monorepo

- **Статус:** Accepted
- **Контекст:** Один разработчик, общие proto/event schemas, проще CI на старте.
- **Решение:** Go monorepo `socset` с `services/`, `pkg/`, `docs/`, `deploy/`.
- **Последствия:** Позже можно вынести сервисы в отдельные репозитории без смены архитектуры.

---

## ADR-010: Observability с первого рабочего скелета

- **Статус:** Accepted
- **Контекст:** Без трейсов микросервисы «непрозрачны».
- **Решение:** OpenTelemetry tracing + structured logs с `request_id` / `trace_id` с фазы 1.
- **Последствия:** Чуть больше boilerplate; проще дебаг интеграций.
