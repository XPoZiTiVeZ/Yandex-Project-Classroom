.PHONY:
	run
	build
	clean
	help

GO=go

help:
	@echo "Доступные команды:"
	@echo "  dev-run  - Запустить проект без сборки"
	@echo "  build    - Собрать проект"
	@echo "  run      - Запустить проект"
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