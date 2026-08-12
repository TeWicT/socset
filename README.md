# Socset

Мини-социальная сеть как pet-проект с микросервисной архитектурой, приближенной к прод-практикам бигтеха.

**Стек:** Go · PostgreSQL · Kafka · Elasticsearch · Redis · MinIO · OpenTelemetry

## Возможности

1. Регистрация и авторизация
2. Профиль с постами на стене, лайками и комментариями; заявки в друзья
3. Поиск профилей (и постов)
4. Страница друзей
5. Чаты по WebSocket
6. Уведомления (in-app / позже push)

## Сервисы

| Сервис | Ответственность |
|--------|-----------------|
| [api-gateway](docs/services/api-gateway.md) | JWT authn, routing, rate limit, request id |
| [auth-service](docs/services/auth-service.md) | Регистрация, логин, токены, credentials |
| [profile-service](docs/services/profile-service.md) | Профиль, настройки, аватар |
| [social-graph-service](docs/services/social-graph-service.md) | Друзья и заявки |
| [post-service](docs/services/post-service.md) | Стена, посты, лайки, комментарии |
| [chat-service](docs/services/chat-service.md) | WebSocket, диалоги, сообщения |
| [notification-service](docs/services/notification-service.md) | In-app / push по событиям |
| [search-service](docs/services/search-service.md) | Индексация в Elasticsearch, поиск |

## Документация

- [Индекс документации](docs/README.md)
- [Архитектура](docs/architecture/overview.md)
- [Каталог событий Kafka](docs/events/catalog.md)
- [Инфраструктура](docs/infra/stack.md)
- [Локальный запуск](docs/local-development.md)
- [Roadmap](docs/roadmap.md)

## Статус

Проект на этапе реализации. Реализация идёт поэтапно — см. [roadmap](docs/roadmap.md).

## Лицензия

MIT (или уточнить позже).
