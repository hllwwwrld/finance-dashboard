include: 

MIGRATIONS_FOLDER:="./db/migrations"
DB_DSN:="postgres://postgres:postgres@localhost:5432/dashboard?sslmode=disable"

bin-deps:
	GOBIN=$(LOCAL_BIN) go install golang.org/x/tools/cmd/goimports@latest
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	GOBIN=$(LOCAL_BIN) go install honnef.co/go/tools/cmd/staticcheck@latest

deps:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest


lint:
	golangci-lint run
	staticcheck ./...

format:
	goimports -w .

# Docker команды используют корневой docker-compose.yml
build:
	docker-compose -f ../docker-compose.yml build --no-cache backend

.PHONY: up
up:
	docker-compose -f ../docker-compose.yml up -d

build-up: build up

.PHONY: down
down:
	docker-compose -f ../docker-compose.yml down

# для локальной разработки
# с удалением volumes в контейнера постегреса, могут потеряться данные
l-up: up
	goose -dir $(MIGRATIONS_FOLDER) postgres $(DB_DSN) up

l-down:
	docker-compose -f ../docker-compose.yml down -v --remove-orphans

.PHONY: restart
restart: down up
