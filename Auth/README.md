# Auth Service

Сервис, отвечающий за регистрацию, авторизацию и аутентификацию. Использует связку access-refresh jwt токенов.

## 📦 Стек

- Go
- gRPC
- PostgreSQL
- Redis
- JWT

## 🚀 Возможности

- Регистрация пользователя
- Вход в аккаунт
- Обновление access токена
- Управление ролями пользователей
- Выход из аккаунта
- Получение информации по пользователю

## ⚙️ Конфигурация

Конфигурация задается через `config.yaml` или env, могут использовать оба способа, env имеют приоритет на yaml

Пример yaml конфигурации:

```yaml
port: 50051
postgres_url: 'postgres://postgres:password@postgres:5432/classroom?sslmode=disable'
redis_url: 'redis://redis:6379/0'
auth:
  jwt_secret: secret
  access_ttl: 15m
  refresh_ttl: 720h
```

Пример env конфигурации:

```bash
PORT=8080
POSTGRES_URL=postgres://postgres:password@localhost:5432/classroom?sslmode=disable
REDIS_URL=redis://localhost:6379/0
AUTH_JWT_SECRET=secret
AUTH_ACCESS_TTL=15m
AUTH_REFRESH_TTL=24h
```

## 🧪 Тестирование

Для написания unit-тестов рекомендуется использовать библиотеку [`mockery`](https://github.com/vektra/mockery) для генерации моков интерфейсов.

### 🔧 Установка mockery

```bash
go install github.com/vektra/mockery/v2@latest
```

### ⚙️ Конфигурация `.mockery.yml`

Добавьте в файл `.mockery.yml` интерфейсы, для которых нужно сгенерировать моки

### 🚀 Генерация моков

После настройки `.mockery.yml` выполните команду:

```bash
mockery
```
