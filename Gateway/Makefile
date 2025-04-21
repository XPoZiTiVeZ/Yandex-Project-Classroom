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
	@echo "  dev-run  - Запустить проект без сборки"
	@echo "  build    - Собрать проект"
	@echo "  run      - Запустить проект"
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