.PHONY:
	run
	build
	clean
	help
	dev-run

GO=go

help:
	@echo "Доступные команды:"
	@echo "  dev-run  - Запустить сервис без сборки"
	@echo "  build    - Собрать сервис"
	@echo "  run      - Запустить сервис"
	@echo "  clean    - Очистить скомпилированные файлы"

dev-run:
	$(GO) run ./cmd

build: 
	$(GO) build -o ./bin/app ./cmd
	
run: build
	./bin/app

clean:
	$(GO) clean
	rm -rf ./bin