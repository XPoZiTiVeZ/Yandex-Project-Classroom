# Опционально загружаем переменные из .env, если файл существует, для POSTGRES_URL
ifneq (,$(wildcard .env))
	include .env
	export
endif

POSTGRES_URL ?= postgres://user:password@localhost:5432/dbname?sslmode=disable

BINARY_NAME=classroom

GO=go

AUTH_GEN_DIR=./Auth/pkg/
COURSE_GEN_DIR=./Courses/pkg/
LESSON_GEN_DIR=./Lessons/pkg/
GATEWAY_GEN_DIR=./Gateway/pkg/
NOTIFICATION_GEN_DIR=./Notifications/pkg/
CHAT_GEN_DIR=./Chat/pkg/

PROTO_DIR=./Common/Proto

MIGRATIONS_DIR=./Common/Migrations

# Создание миграций, используется https://github.com/golang-migrate/migrate, name - название миграции
migrate-create:
	@name=$(name);
	@migrate create -seq -ext sql -dir $(MIGRATIONS_DIR) $(name)

# Применение миграций
migrate-up:
	@migrate -path=$(MIGRATIONS_DIR) -database=$(POSTGRES_URL) up

# Откат миграций, name - к какой миграции откатить
migrate-down:
	@name=$(name);
	@migrate -path=$(MIGRATIONS_DIR) -database=$(POSTGRES_URL) down $(name)

# @protoc --proto_path=common/proto/$(AUTH_GEN_DIR) \
# --go_out=common/api/$(name) --go_opt=paths=source_relative \
# --go-grpc_out=common/api/$(name) --go-grpc_opt=paths=source_relative \
# common/api/$(name)/$(name).proto

# @mkdir -p $(NOTIFICATION_GEN_DIR)
# @mkdir -p $(CHAT_GEN_DIR)

## generate: Генерация кода из .proto файлов
proto-gen:
	@mkdir -p $(AUTH_GEN_DIR)
	@mkdir -p $(COURSE_GEN_DIR)
	@mkdir -p $(LESSON_GEN_DIR)
	@mkdir -p $(GATEWAY_GEN_DIR)
	
	@protoc --go_out=$(AUTH_GEN_DIR) --go-grpc_out=$(AUTH_GEN_DIR)  -I. $(PROTO_DIR)/auth.proto
	@protoc --go_out=$(COURSE_GEN_DIR) --go-grpc_out=$(COURSE_GEN_DIR)  -I. $(PROTO_DIR)/course.proto
	@protoc --go_out=$(LESSON_GEN_DIR) --go-grpc_out=$(LESSON_GEN_DIR)  -I. $(PROTO_DIR)/lesson.proto $(PROTO_DIR)/task.proto

	@protoc --go_out=$(GATEWAY_GEN_DIR) --go-grpc_out=$(GATEWAY_GEN_DIR)  -I. $(PROTO_DIR)/auth.proto $(PROTO_DIR)/course.proto $(PROTO_DIR)/lesson.proto $(PROTO_DIR)/task.proto
	

build:
	if [ ! -f "go.mod" ]; then \
		go mod init app; \
	fi
	$(GO) mod tidy
	$(GO) build -o ./bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

## run: Запустить проект
run:
	./bin/$(BINARY_NAME)

## clean: Очистить скомпилированные файлы
clean:
	$(GOCLEAN)
	rm -f ./bin/$(BINARY_NAME)

## help: Показать доступные команды
help:
	@echo "Доступные команды:"
	@echo "  generate - Сгенерировать go файлы из proto\n\
             !!! должен быть установлен protobuf-compiler"
	@echo "  build    - Собрать проект"
	@echo "  run      - Запустить проект"
	@echo "  clean    - Очистить скомпилированные файлы"
	@echo "  migrate-create name=<название миграции> - Создание миграций"
	@echo "  migrate-up - Применение миграций"
	@echo "  migrate-down name=<к какой миграции откатить> - Откат миграций"
	@echo "  help - Показать доступные команды"

# По умолчанию запускается help
.DEFAULT_GOAL := help