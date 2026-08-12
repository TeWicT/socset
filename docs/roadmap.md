# Roadmap

Реализация поэтапная. Каждая фаза даёт работающий вертикальный срез.

## -  Фаза 0 — Документация и каркас репозитория

**Статус:** Готова

- [x] Архитектура, ADR, владение данными
- [x] Описание сервисов
- [x] Каталог событий
- [x] Инфра и локальная разработка
- [x] Monorepo skeleton (`services/`, `pkg/`, `deploy/`)
- [x] `docker-compose` с инфраструктурой
- [x] `.env.example`, Makefile/Task

**Критерий готовности:** можно поднять Postgres/Kafka/Redis/ES/MinIO/Jaeger одной командой.

---

## Фаза 1 — Identity skeleton

**Статус:** В процессе

**Сервисы:** api-gateway, auth-service, profile-service (минимум)

- [x] Регистрация / логин / refresh / logout
- [-] JWT на gateway
- [-] Создание профиля по `user.registered` (outbox → Kafka)
- [-] `GET /profiles/{id}`, `PATCH /profiles/me`
- [-] OTel + structured logs
- [-] OpenAPI для auth и profile

**Критерий:** пользователь регистрируется, логинится, видит/редактирует профиль через gateway.

---

## Фаза 2 — Социальный граф и стена

**Сервисы:** social-graph-service, post-service

- Заявки в друзья, список друзей, status
- Посты на своей стене, лайки, комментарии
- События `friend.*`, `post.*`, `comment.*`
- Аватар через MinIO

**Критерий:** два пользователя дружатся; на стене есть пост с лайком и комментарием.

---

## Фаза 3 — Поиск и уведомления

**Сервисы:** search-service, notification-service

- Индексация users/posts в Elasticsearch
- `GET /search/users`, `GET /search/posts`
- In-app уведомления из Kafka
- Unread count + mark read

**Критерий:** поиск находит пользователя по имени; лайк создаёт уведомление.

---

## Фаза 4 — Чаты

**Сервисы:** chat-service

- DM диалоги, история
- WebSocket real-time
- Redis presence + pub/sub между инстансами
- `message.sent` → notification

**Критерий:** два онлайн-пользователя переписываются в реальном времени.

---

## Фаза 5 — Прод-гигиена

- gRPC между сервисами (где выгодно)
- Idempotency keys на критичных POST
- Нагрузочные smoke-тесты
- CI: lint, test, build images
- Политики ретеншна Kafka, бэкапы Postgres (описание)
- (Опц.) Kubernetes manifests / Kind

**Критерий:** CI зелёный; README описывает полный локальный сценарий E2E.

---

## Вне скоупа MVP (бэклог)

- Групповые чаты, вложения в сообщениях
- Лента новостей друзей
- Push (FCM/WebPush)
- Модерация / жалобы
- Рекомендации «люди, которых вы можете знать»
- Multi-region, CQRS-лента, GraphQL BFF

## Текущий фокус

После документации — **Фаза 0 (каркас) → Фаза 1**.
