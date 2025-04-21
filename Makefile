.PHONY:
	migrate-create
	migrate-up
	migrate-down
	proto-gen
	proto-gen-gateway
	proto-gen-auth
	proto-gen-course
	proto-gen-lessons
	proto-gen-tasks
	proto-gen-notifications
	proto-gen-chat
	generate-doc
	build
	run
	clean
	help

# Опционально загружаем переменные из .env, если файл существует, для POSTGRES_URL
ifneq (,$(wildcard .env))
	include .env
	export
endif

POSTGRES_URL ?= postgres://user:password@localhost:5432/dbname?sslmode=disable

BINARY_NAME=classroom

GO=go

GATEWAY_GEN_DIR=./Gateway/pkg/

AUTH_GEN_DIR=./Auth/pkg/
COURSES_GEN_DIR=./Courses/pkg/
LESSONS_GEN_DIR=./Lessons/pkg/
TASKS_GEN_DIR=./Tasks/pkg/
NOTIFICATIONS_GEN_DIR=./Notifications/pkg/
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

## generate: Генерация кода из .proto файлов
proto-gen: \
	proto-gen-gateway \
	proto-gen-auth \
	proto-gen-courses \
	proto-gen-lessons \
	proto-gen-tasks \
# proto-gen-notifications \
proto-gen-chat \

proto-gen-gateway:
	@mkdir -p $(GATEWAY_GEN_DIR)
	@protoc --go_out=$(GATEWAY_GEN_DIR) --go-grpc_out=$(GATEWAY_GEN_DIR)  -I. $(wildcard $(PROTO_DIR)/*.proto)

proto-gen-auth:
	@mkdir -p $(AUTH_GEN_DIR)
	@protoc --go_out=$(AUTH_GEN_DIR) --go-grpc_out=$(AUTH_GEN_DIR)  -I. $(PROTO_DIR)/auth.proto

proto-gen-courses:
	@mkdir -p $(COURSES_GEN_DIR)
	@protoc --go_out=$(COURSES_GEN_DIR) --go-grpc_out=$(COURSES_GEN_DIR)  -I. $(PROTO_DIR)/courses.proto

proto-gen-lessons:
	@mkdir -p $(LESSONS_GEN_DIR)
	@protoc --go_out=$(LESSONS_GEN_DIR) --go-grpc_out=$(LESSONS_GEN_DIR)  -I. $(PROTO_DIR)/lessons.proto

proto-gen-tasks:
	@mkdir -p $(TASKS_GEN_DIR)
	@protoc --go_out=$(TASKS_GEN_DIR) --go-grpc_out=$(TASKS_GEN_DIR)  -I. $(PROTO_DIR)/tasks.proto

proto-gen-notifications:
	@mkdir -p $(NOTIFICATIONS_GEN_DIR)
	@protoc --go_out=$(NOTIFICATIONS_GEN_DIR) --go-grpc_out=$(NOTIFICATIONS_GEN_DIR)  -I. $(PROTO_DIR)/notifications.proto

proto-gen-chat:
	@mkdir -p $(CHAT_GEN_DIR)
	@protoc --go_out=$(CHAT_GEN_DIR) --go-grpc_out=$(CHAT_GEN_DIR)  -I. $(PROTO_DIR)/chat.proto

generate-doc:
	cd Gateway/ && \
	swag init -g ./cmd/main.go --output ./docs --parseDependency --parseInternal --outputTypes json && \
	cd ..

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
	@echo "  proto-gen - Сгенерировать go файлы из proto для всех сервисов\n\
              !!! должен быть установлен protobuf-compiler"
	@echo "  proto-gen-<сервис> - Сгенерировать go файлы из proto для <сервис>\n\
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