# API-Gateway Service

Сервис для обработки приходящих http запросов. Реализует доступ ко всем функциям остальных сервисов, кроме служебных.

## 📦 Стек

- Go
- gRPC
- Redis
- Kafka

## 📃 Возможности

- Приём http запросов и перенаправление их на grpc микросервисы.
- Сбор необходимой информации запросу с микросервисов
- Авторизация пользователей
- Отправка уведомлений через сервис уведомлений

## ⚙️ Конфигурация

Конфигурация задается через `config.yaml` в папке `config/` сервиса и переменные окружения:

```yaml
host:
  address: 0.0.0.0
  port: 8080
auth:
  address: service-auth
  port: 50051
  enabled: false
courses:
  address: service-courses
  port: 50052
  enabled: false
lessons:
  address: service-lessons
  port: 50053
  enabled: false
tasks:
  address: service-tasks
  port: 50054
  enabled: false
notifications:
  address: service-notifications
  port: 50055
  enabled: false
chat:
  address: service-chat
  port: 50056
  enabled: false
```

В config.yaml прописываются:

- адрес, на котором работает сервис.
- адреса используемых сервисов, а также их статус включения, в зависимости их статуса готовности.

В переменных окружения задаются чувствительные данные:

```env
REDIS_URL=redis://redis:6379/0
KAFKA_URL=kafka://kafka:9092
AUTH_JWT_SECRET=secret
```

Аргумент командной строки для конфига:

```
-config ./config/config.yaml
```

## 🧪 Тестирование

- Юнит-тестирование сервиса невозможно ввиду его функций.
- Интеграционные тесты пока не написаны.

<!-- ### 🔧 Установка mockery

```bash
go install github.com/vektra/mockery/v2@latest
```

### ⚙️ Конфигурация `.mockery.yml`

Добавьте в файл `.mockery.yml` интерфейсы, для которых нужно сгенерировать моки

### 🚀 Генерация моков

После настройки `.mockery.yml` выполните команду:

```bash
mockery
``` -->
