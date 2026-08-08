# Документация Socset

Точка входа в проектную документацию.

## Разделы

| Раздел | Описание |
|--------|----------|
| [Архитектура](architecture/overview.md) | C4-уровень, потоки запросов, принципы |
| [Владение данными](architecture/data-ownership.md) | Какой сервис чем владеет |
| [ADR](architecture/decisions.md) | Architecture Decision Records |
| [Сервисы](services/) | Контракты и границы каждого сервиса |
| [События Kafka](events/catalog.md) | Топики, схемы, consumer'ы |
| [Инфраструктура](infra/stack.md) | Postgres, Kafka, ES, Redis, MinIO, observability |
| [Локальная разработка](local-development.md) | docker-compose, порядок подъёма |
| [Roadmap](roadmap.md) | Фазы реализации |
| [API](api/README.md) | Внешние HTTP/WS контракты |

## Как читать

1. Начни с [overview](architecture/overview.md) и [roadmap](roadmap.md).
2. Перед реализацией сервиса — его файл в `services/` + связанные события.
3. Перед интеграцией — [catalog](events/catalog.md) и [data-ownership](architecture/data-ownership.md).

## Соглашения

- Документация на русском; идентификаторы (сервисы, топики, поля) — на английском.
- Изменения архитектуры фиксируются ADR в `architecture/decisions.md`.
- Публичные API описываются в OpenAPI (по мере появления кода) и кратко дублируются в `docs/api/`.
