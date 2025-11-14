# testTask – сервис вопросов и ответов

REST-сервис на Go для работы с вопросами и ответами.  

## Стек

- Go (Golang)
- PostgreSQL
- GORM
- goose (SQL-миграции)
- Docker / docker-compose

---

## Подготовка `.env`

В корне проекта создайте файл `.env` (если его ещё нет) и заполните, например так:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=testTask
DB_PORT=5432
DB_SSLMODE=disable
APP_PORT=:8080
DSN_PG="postgres://postgres:secret@db:5432/testTask?sslmode=disable"
```

## Запуск

В папке проекта собрать и поднять все сервисы:
```bash
docker-compose up --build
```
Для удобства в корне проекта лежит коллекция в Insomnia:

Insomnia_2025-11-14.yaml

## Краткое описание API

Проект реализует простый CRUD для вопросов и ответов (Questions / Answers):

создание вопроса;

получение списка вопросов;

добавление ответа к вопросу;

получение/удаление ответа и вопроса.
