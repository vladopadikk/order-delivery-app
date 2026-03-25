# Order Delivery App

Микросервисное backend-приложение для регистрации пользователя, оформления, оплаты, доставки и уведомления о заказах.

Проект демонстрирует микросервисную архитектуру, проектирование REST API, использование PostgreSQL, Kafka, Docker, 
применение JWT-аутентификации, транзакций, миграций и unit-тестирование.

---

## 📌 О проекте

Система моделирует полный цикл обработки заказа:

- пользователь регистрируется и авторизуется;
- создает заказ;
- заказ публикуется как событие в Kafka;
- сервис оплаты обрабатывает оплату;
- сервис доставки создает и завершает доставку;
- сервис уведомлений отправляет пользователю сообщение с информацией о состоянии заказа.

Проект построен в микросервисном стиле:
- каждый сервис отвечает за свою бизнес-область;
- сервисы общаются через Kafka;
- каждый сервис с БД использует собственную PostgreSQL базу данных.

---

## 🧱 Архитектура

Сервисы:
- `auth-service` — регистрация, логин, JWT
- `orders-service` — создание и просмотр заказов
- `payment-service` — обработка оплаты заказа
- `delivery-service` — создание и завершение доставки
- `notification-service` — обработка событий и отправка уведомлений


Взаимодействие:
- **REST API** — клиентские запросы
- **Kafka** — асинхронный обмен событиями между сервисами

---

## ⚙️ Технологии

- **Go**
- **Gin**
- **PostgreSQL**
- **Kafka**
- **Docker / Docker Compose**
- **JWT**
- **Goose migrations**
- **Unit tests**
- **Git**

---

## 🔄 Основной сценарий работы

1. Пользователь регистрируется / логинится через `auth-service`
2. Получает JWT-токен
3. Создает заказ через `orders-service`
4. `orders-service` сохраняет заказ и публикует событие `order_created`
5. `payment-service` получает событие и имитирует оплату
6. В зависимости от результата публикуется:
   - `payment_success`
   - `payment_failed`
7. `delivery-service` имитирует доставку после успешной оплаты
8. Публикуется событие `delivery_completed`
9. `orders-service` получает событие и изменяет состояние заказа 
10. `notification-service` отправляет уведомление пользователю

---

## 📨 Kafka Events

Используемые топики:

- `order_created`
- `payment_success`
- `payment_failed`
- `delivery_completed`

Пример event-driven цепочки:

`orders-service` → `order_created` → `payment-service`  

`payment-service` → `payment_success` → `delivery-service`  

`payment-service` → `payment_success` / `payment_failed` → `orders-service`

`payment-service` → `payment_success` / `payment_failed` → `notification-service`

`delivery-service` → `delivery_completed` → `orders-service`

`delivery-service` → `delivery_completed` → `notification-service`

---

## 🗂️ Структура проекта

```text
├── auth-service
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── config
│   │   │   └── config.go
│   │   ├── database
│   │   │   └── connect.go
│   │   ├── handlers
│   │   │   └── handler.go
│   │   ├── models
│   │   │   ├── login.go
│   │   │   └── user.go
│   │   ├── repository
│   │   │   └── repo.go
│   │   └── service
│   │       ├── service.go
│   │       └── tokens.go
│   ├── main.go
│   └── migrations
│       └── 20260207121707_create_users_table.sql
├── delivery-service
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── config
│   │   │   └── config.go
│   │   ├── database
│   │   │   └── connect.go
│   │   ├── handlers
│   │   │   └── handler.go
│   │   ├── kafka
│   │   │   ├── consumer
│   │   │   │   └── consumer.go
│   │   │   └── producer
│   │   │       └── producer.go
│   │   ├── middleware
│   │   │   └── auth.go
│   │   ├── models
│   │   │   ├── deliveryDTO.go
│   │   │   ├── delivery.go
│   │   │   └── events.go
│   │   ├── repository
│   │   │   └── repo.go
│   │   └── service
│   │       └── service.go
│   ├── main.go
│   └── migrations
│       └── 20260215091046_create_delivery_table.sql
├── docker-compose.yml
├── go.mod
├── go.sum
├── kafka
│   └── initTopics.go
├── notification-service
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── config
│   │   │   └── config.go
│   │   ├── consumer
│   │   │   └── consumer.go
│   │   ├── models
│   │   │   └── events.go
│   │   └── service
│   │       ├── service.go
│   │       └── service_test.go
│   └── main.go
├── orders-service
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── config
│   │   │   └── config.go
│   │   ├── database
│   │   │   ├── connect.go
│   │   │   └── executor.go
│   │   ├── handlers
│   │   │   └── handler.go
│   │   ├── kafka
│   │   │   ├── consumer
│   │   │   │   └── consumer.go
│   │   │   └── producer
│   │   │       └── producer.go
│   │   ├── middleware
│   │   │   └── auth.go
│   │   ├── models
│   │   │   ├── events.go
│   │   │   ├── orderDTO.go
│   │   │   └── order.go
│   │   ├── repository
│   │   │   └── repo.go
│   │   └── service
│   │       └── service.go
│   ├── main.go
│   └── migrations
│       └── 20260208090014_cretate_orders_table.sql
├── payments-service
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── config
│   │   │   └── config.go
│   │   ├── database
│   │   │   └── connect.go
│   │   ├── handlers
│   │   │   └── handler.go
│   │   ├── kafka
│   │   │   ├── consumer
│   │   │   │   └── consumer.go
│   │   │   └── producer
│   │   │       └── producer.go
│   │   ├── middleware
│   │   │   └── auth.go
│   │   ├── models
│   │   │   ├── events.go
│   │   │   ├── payment.go
│   │   │   └── paymentsDTO.go
│   │   ├── repository
│   │   │   └── repo.go
│   │   └── service
│   │       └── service.go
│   ├── main.go
│   └── migrations
│       └── 20260211082611_create_payments_table.sql
└── README.md
```

---

## 🧠 Ключевые технические решения
- Микросервисная архитектура — каждый сервис отвечает за отдельную бизнес-область
- Отдельная БД на сервис — изоляция данных и независимость сервисов
- Kafka для асинхронного взаимодействия — реализация event-driven flow
- Транзакции в orders-service — создание заказа и его позиций выполняется атомарно
- JWT-аутентификация — защита приватных маршрутов
- Слоистая архитектура:
  handler → service → repository
- Интерфейсы для зависимостей — упрощают unit-тестирование business logic

---

## 🐳 Запуск проекта

### Требования

- Docker Dekstop
- Docker Compose
- goose

### Запуск проекта

1. Поднять zookeeper и kafka
   
```bash
docker compose up -d zookeeper kafka
```

2. Создать топики

```bash
docker exec -it kafka kafka-topics --create --if-not-exists --topic order_created --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1
docker exec -it kafka kafka-topics --create --if-not-exists --topic payment_success --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1
docker exec -it kafka kafka-topics --create --if-not-exists --topic payment_failed --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1
docker exec -it kafka kafka-topics --create --if-not-exists --topic delivery_completed --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1
```

3. Поднять сервисы и БД

```bash
docker compose up -d auth-db orders-db payment-db delivery-db auth-service orders-service payment-service delivery-service notification-service 
```

---

##  REST API эндпоинты

### Публичные эндпоинты

#### Регистрация пользователя

```http
POST /api/register
Content-Type: application/json

{
  "username": "ivan",
  "email": "ivan@example.com",
  "password": "password123"
}
```

**Ответ:** `201 Created`
```json
{
  "id": 1,
  "username": "ivan",
  "email": "ivan@example.com"
}
```

**Ошибки:**
- `400 Bad Request` - invalid JSON
- `409 Conflict` - email already exists

---

#### Логин
```http
POST /api/login
Content-Type: application/json

{
  "email": "ivan@example.com",
  "password": "password123"
}
```

**Ответ:** `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Ошибки:**
- `400 Bad Request` - invalid JSON
- `404 Not Found` - user not found
- `401 Unauthorized` - invalid password

---

### Защищённые эндпоинты
Все эндпоинты ниже требуют заголовок: `Authorization`:

```
Authorization: Bearer <access_token>
```

---

#### Создание заказа
```http
POST /orders
Content-Type: application/json
Authorization: Bearer <token>

{
  "items": [
    {
      "item_id": 12345,
      "quantity": 2,
      "price": 29.99
    },
    {
      "item_id": 67890,
      "quantity": 1,
      "price": 149.50
    }
  ],
  "delivery_address": "г. Москва, ул. Тверская, д. 10, кв. 25"
}
```

---

#### Получение списка заказов
```http
GET /orders
Content-Type: application/json
Authorization: Bearer <token>
```

---

#### Получение статуса доставки заказа
```http
GET /deliveries?order_id=1
Content-Type: application/json
Authorization: Bearer <token>
```

**Query Parameters:**
- `order_id` - идентификатор заказа

---

#### Получение статуса оплаты заказа
```http
GET /payments?order_id=1
Content-Type: application/json
Authorization: Bearer <token>
```

**Query Parameters:**
- `order_id` - идентификатор заказа

---

## ✅ Тестирование

На текущем этапе в проекте реализованы базовые unit-тесты для бизнес-логики `notification-service`.

Тесты покрывают **service-слой** и проверяют, что входящие события корректно преобразуются в пользовательские уведомления.  
Для изоляции бизнес-логики от внешней инфраструктуры используется mock-реализация `Notifier`.

Покрытые сценарии:
- успешная оплата заказа
- неуспешная оплата заказа
- успешное завершение доставки

Что проверяется в тестах:
- вызывается notifier
- в notifier передаётся корректный `userID`
- формируется корректный текст уведомления для каждого бизнес-сценария

Такой подход демонстрирует:
- использование интерфейсов для внедрения зависимостей
- разделение бизнес-логики и инфраструктурного слоя
- умение писать unit-тесты без запуска Kafka, HTTP-обработчиков и базы данных

---
