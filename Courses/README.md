# Courses Service

Сервис, отвечающий за курсы. Позоляет создавать, получать, редактировать, удалять курсы и управлять студентами курса

## 📦 Стек

- Go
- gRPC
- PostgreSQL

## 🚀 Возможности

- Создание курсов
- Управление курсами
- Добавление новых пользователей на курс
- Получение информации об курсах
- Получение списка студентов курса
- Получение списка курсов пользователя 

## ⚙️ Конфигурация

Конфигурация задается через `config.yaml` или env, могут использовать оба способа, env имеют приоритет на yaml

Пример yaml конфигурации:

```yaml
port: 50051
postgres_url: 'postgres://postgres:password@courses:5432/classroom?sslmode=disable'

```

Пример env конфигурации:

```bash
PORT=50052
POSTGRES_URL=postgres://postgres:password@localhost:5432/classroom?sslmode=disable
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
