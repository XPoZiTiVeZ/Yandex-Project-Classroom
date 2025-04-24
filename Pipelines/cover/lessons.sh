#!/bin/bash

# * Папки для исключения из покрытия
EXCLUDE_DIRS=(
  "mocks"
  "pkg"
  "cmd"
  "internal/config"
  "internal/dto"
  "internal/domain"
  "internal/producer"
  "internal/repo"
  "internal/service/mocks"
  "internal/controller/mocks"
)

# * Папка с микросевисом 
# ! Обязательно изменить на свой
MICROSERVICE_NAME="Lessons"

cd "$MICROSERVICE_NAME"
EXCLUDE_PATHS=$(IFS="|"; echo "${EXCLUDE_DIRS[*]}")
PATHS_FOR_TESTS=$(go list ./... | grep -v -E "$EXCLUDE_PATHS")

# * Тестирование
# * Если нужна доп. информация, закомментируйте '> /dev/null 2>&1'
go test -coverprofile=coverage.out $PATHS_FOR_TESTS > /dev/null 2>&1
COVERAGE_PERCENT=$(go tool cover -func=coverage.out | grep "total:" | awk '{print int($3)}')
rm coverage.out

cd ..

echo "$COVERAGE_PERCENT"
