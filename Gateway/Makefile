.PHONY:
	run
	build
	clean
	help

GO=go
SWAG=swag

help:
	@echo "Доступные команды:"
	@echo "  generate-doc - Сгенерировать документацию API"
	@echo "  dev-run  - Запустить сервис без сборки"
	@echo "  build    - Собрать сервис"
	@echo "  run      - Запустить сервис"
	@echo "  clean    - Очистить скомпилированные файлы"

generate-doc:
	$(SWAG) init -g ./cmd/main.go --output ./docs --parseDependency --parseInternal --outputTypes json

dev-run:
	$(GO) run ./cmd

build: 
	$(GO) build -o ./bin/app ./cmd
	
run: build
	./bin/app

clean:
	$(GO) clean
	rm -rf ./bin