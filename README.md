# HireRadar

Автоматизированная система поиска и фильтрации вакансий с hh.ru с многоцелевым таргетингом и доставкой уведомлений через Telegram Bot и HTTP API.

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-4169E1?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7+-DC382D?style=flat&logo=redis)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-✓-2496ED?style=flat&logo=docker)](https://www.docker.com/)


## О проекте

Поиск релевантных кандидатов или вакансий на hh.ru вручную отнимает часы у HR-специалистов и соискателей. HireRadar решает эту задачу: система по заданным фильтрам автоматически собирает вакансии, дедуплицирует их, сохраняет в собственную базу и доставляет пользователю через удобный Telegram-бот или HTTP API. Планировщик запускает фоновый сбор по расписанию, а слой кэширования убирает дубли и снижает нагрузку на внешний API.

## Ключевые возможности

- **Многоканальная доставка** — Telegram Bot и HTTP API для интеграции c внешними сервисами.
- **Умные фильтры** — поиск по названию, городу, зарплатным ожиданиям и другим параметрам.
- **Дедупликация вакансий** — Redis-слой гарантирует, что пользователь не увидит одну и ту же вакансию дважды.
- **Поисковые агенты по расписанию** — встроенный планировщик с настраиваемым интервалом опроса hh.ru.
- **Персональные аккаунты** — JWT-аутентификация, избранное и сохранённые поисковые профили.
- **Миграции БД** — версионирование схемы PostgreSQL из коробки.

## Стек

| Компонент          | Технология                                        |
| -------------------- | --------------------------------------------------- |
| Язык               | Go 1.23+                                          |
| База данных        | PostgreSQL 15+ (драйвер `pgx`)                    |
| Кэш / дедупликация | Redis 7+                                          |
| Транспорт          | Telegram Bot API, REST (Gin)                      |
| Аутентификация     | JWT                                               |
| Планировщик        | Встроенный cron-подобный модуль                   |
| Миграции           | `golang-migrate`                                  |
| Тестирование       | `testify`, `testcontainers-go`, in-memory mock'и |
| Контейнеризация    | Docker, Docker Compose                            |

## Архитектура


1. **Planner** (`internal/scheduler`) по расписанию вызывает поиск через hh.ru API.
2. Полученные вакансии проходят через доменную логику фильтрации и сохраняются в PostgreSQL.
3. Redis используется как быстрый слой проверки «уже видели», чтобы исключить дубли.
4. Пользователи взаимодействуют с системой либо через Telegram-бота, либо напрямую через REST API.

## Быстрый старт

### Предварительные требования

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- (Опционально) Go 1.23+ для локальной разработки

### Запуск через Docker Compose

```bash
# 1. Клонируйте репозиторий
git clone https://github.com/your-org/hireradar.git
cd hireradar

# 2. Создайте .env файл (пример ниже)
cp .env.example .env

# 3. Поднять только инфраструктуру (БД + Redis)
docker compose up -d postgres redis

# 4. Установить утилиту migrate (если ещё нет)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 5. Применить миграции
make migrate-up

# 6. Запустить приложение
go run cmd/bot/main.go
```

## API Reference
```
API Reference
Базовый URL: http://localhost:8080/api/v1

Метод	Эндпоинт	Защита	Назначение
POST	/auth/register	Нет	Регистрация нового пользователя
POST	/auth/login	Нет	Получение JWT-токена
GET	/vacancies	JWT	Список вакансий, отфильтрованных под профиль
POST	/filters	JWT	Создать поисковый фильтр
GET	/filters	JWT	Получить свои фильтры
POST	/favorites	JWT	Добавить вакансию в избранное
GET	/favorites	JWT	Список избранного
```
## Структура проекта

```
.

├── cmd/bot/                # Точка входа
├── internal/
│   ├── client/hh/          # HTTP-клиент для API hh.ru
│   ├── config/             # Загрузка конфигурации
│   ├── delivery/
│   │   ├── bot/            # Telegram Bot (telebot)
│   │   └── http/           # REST API (Gin)
│   ├── domain/             # Бизнес-сущности и интерфейсы
│   ├── mocks/              # Mock'и для тестов
│   ├── repository/
│   │   ├── benchmarks/     # Бенчмарки in-memory, Postgres, Redis
│   │   ├── postgres/       # Реализации репозиториев на pgx
│   │   └── redis/          # Кэш и дедупликация
│   ├── scheduler/          # Планировщик фонового сбора
│   ├── usecase/            # Слой бизнес-логики
│   └── utils/              # Миграции и хелперы
├── migrations/             # SQL-миграции
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── go.mod
```