# Balance Service

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue.svg)](https://go.dev/)

RESTful сервис на Go для управления балансом пользователей. Сервис реализует базовые финансовые операции: начисление средств, получение баланса, резервирование средств и создание финансовой отчётности. Проект построен с акцентом на чистоту архитектуры и консистентность данных.

## Стек технологий
- Язык: [Go](https://go.dev/) 
- Стандартная библиотека: [net/http](https://pkg.go.dev/net/http) 
- База данных: [PostgreSQL](https://www.postgresql.org/) 
- Управление конфигурацией: [flag](https://pkg.go.dev/flag), [godotenv](https://github.com/joho/godotenv) 
- Миграции БД: [goose](https://github.com/pressly/goose)

## Установка и запуск

### 1. Клонирование репозитория
`git clone https://github.com/YAMI4YOU/balance-service.git` \
`cd balance-service`

### 2. Настройка окружения
Скопируйте файл `.env.example` в `.env` и укажите свои данные для подключения к БД. \
`cp .env.example .env`

### 3. Настройка базы данных
Установите goose. \
`go install github.com/pressly/goose/v3/cmd/goose@latest` \
Примените миграции для создания необходимых таблиц. \
`goose up`

### 4. Запуск сервиса
`go run ./cmd/main.go`

## Использование флагов командной строки

- Запуск сервиса на порту 3000 \
`go run ./cmd/main.go -port=3000`
- Запуск с указанием адреса БД \
`go run ./cmd/main.go -dburl="postgres://user:pass@host:5432/dbname"`

## API эндпоинты
- Получение средств - `/balance` 
```
curl http://localhost:8080/balance?user_id=1

```
- Начисление средств - `/deposit`
```
curl -X POST http://localhost:8080/deposit \
-H 'Content-Type: application/json' \
-d '{
    "user_id": 2,
    "balance": 12.54
}'
```
- Резервирование средств - `/reserve`
```
curl -X POST http://localhost:8080/reserve \
-H 'Content-Type: application/json' \
-d '{
    "user_id": 1,
    "service_id": 2,
    "order_id": 3,
    "amount": 5000
}'
```
- Подтверждение выручки - `/revenue`
```
curl -X POST http://localhost:8080/revenue \
-H 'Content-Type: application/json' \
-d '{
    "user_id": 1,
    "service_id": 2,
    "order_id": 3,
    "amount": 5000
}'
```
- Формирование отчёта - `/report`
```
curl -X GET http://localhost:8080/report \
-H 'Content-Type: application/json' \
-d '{
    "year": 2025,
    "month": 9
}'
```