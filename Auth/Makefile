include .env
MIGRATIONS_PATH=./migrations

migrate-create:
	@name=$(name);
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(name)

migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(POSTGRES_URL) up

migrate-down:
	@name=$(name);
	@migrate -path=$(MIGRATIONS_PATH) -database=$(POSTGRES_URL) down $(name)

gen-proto:
	@protoc --proto_path=pkg/api/auth \
	  --go_out=pkg/api/auth --go_opt=paths=source_relative \
		--go-grpc_out=pkg/api/auth --go-grpc_opt=paths=source_relative pkg/api/auth/auth.proto
		
run:
	@go run cmd/main.go

test:
	@go test ./... -cover
