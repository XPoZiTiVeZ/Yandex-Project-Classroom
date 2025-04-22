

# Classroom Backend

`classroom` — это сервис, предоставляющий функциональность для управления курсами, пользователями, домашними заданиями и взаимодействием между участниками образовательного процесса.

## 🛠️ Стек технологий

- **Go** — основной язык разработки
- **gRPC** — взаимодействие между модулями и внешними сервисами
- **PostgreSQL** — основное хранилище данных
- **Redis** — кэширование и хранение сессий
- **Kafka** — асинхронная обработка событий
- **JWT** — аутентификация и авторизация
- **Mockery** — генерация моков для тестирования

## 📦 Основные возможности

- Создание и редактирование курсов
- Добавление и удаление студентов
- Управление заданиями и домашними работами
- Аутентификация через JWT
- Уведомление пользователей вместе с помощью Kafka
- Кэширование данных и хранение сессий в Redis

## 🏗️ Архитектура

- **Монорепозитарная структура** - Структура проекта, в которой каждый сервис в своей папке, а не в каждом репозитории.
- **Микросервисная архитектура** - Архитектура проекта строится на микросервисах, каждый отвечает строго за свои задачи.
- gRPC используется как внутренний транспортный протокол между слоями.
- Kafka отвечает за хранение заданий на отсылку уведомлений сервису уведомлений.
- Redis ускоряет доступ к часто запрашиваемым данным (например, профили пользователей, активные сессии).
- PostgreSQL хранит все основные бизнес-данные.

## 🏢 Сервисы

- [Микросервис аутентификации](./Auth/)
- [Микросервис управления курсами](./Courses/)
- [Микросервис управления уроками](./Lessons/)
- [Микросервис управления заданиями](./Tasks/)
- [Микросервис уведомлений](./Notifications/)

## 📃 Настройка проекта

### **Клонируйте репозиторий**

```bash
git clone https://gitlab.crja72.ru/golang/2025/spring/course/projects/go1/classroom.git
cd classroom
```

### **Создайте файл `.env`**

Убедитесь, что в файле `.env` содержатся все необходимые переменные:

```env
# Nginx
NGINX_PORT=80


# Postgres
POSTGRES_DB=classroom
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_URL="postgres://postgres:password@postgres:5432/classroom?sslmode=disable"


# Redis
REDIS_PORT=6379
# REDIS_PASSWORD=1324
REDIS_URL="redis://redis:6379/0"


# Kafka
KAFKA_PORT=9092
KAFKA_URL="kafka://broker1:9092"


# Auth
AUTH_PORT=50051
AUTH_ADDRESS=auth
AUTH_ACCESS_TTL=15m
AUTH_REFRESH_TTL=24h
AUTH_JWT_SECRET=secret
AUTH_ENABLED=true


# Courses
COURSES_PORT=50052
COURSES_ADDRESS=courses
COURSES_ENABLED=true


# Lessons
LESSONS_PORT=50053
LESSONS_ADDRESS=lessons
LESSONS_ENABLED=true


# Tasks
TASKS_PORT=50054
TASKS_ADDRESS=tasks
TASKS_ENABLED=true


# Chat
CHAT_PORT=50055
CHAT_ADDRESS=chat
CHAT_ENABLED=false


# Notifications
NOTIFICATION_PORT=50056
NOTIFICATIONS_ENABLED=fasle


# Gateway
GATEWAY_PORT=8080
GATEWAY_ADDRESS=gateway


# Common
TIMEOUT=10s
MAX_RETRIES=5
```

## 🚀 **Запуск проекта**

### **Запустить проект**

```bash
docker-compose up
```

Если переменных окружения в `.env` достаточно, сервисы запустятся без ошибок.

### **Выключить проект**

```bash
docker-compose down
```

### **Пересобрать и запустить заново**

```bash
docker-compose down
docker-compose up --build
```
